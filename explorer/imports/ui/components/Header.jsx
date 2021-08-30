import qs from 'querystring';
import React,{ Component } from 'react';
import { HTTP } from 'meteor/http'
import {
    Badge,
    Button,
    Collapse,
    Navbar,
    NavbarToggler,
    NavbarBrand,
    Nav,
    NavItem,
    NavLink,
    // Input,
    // InputGroup,
    // InputGroupAddon,
    // Button,
    UncontrolledDropdown,
    UncontrolledPopover,
    PopoverBody,
    DropdownToggle,
    DropdownMenu,
    DropdownItem
} from 'reactstrap';
import { Link } from 'react-router-dom';
import SearchBar from './SearchBar.jsx';
import i18n from 'meteor/universe:i18n';
import LedgerModal from '../ledger/LedgerModal.jsx';
import Account from './Account.jsx';

const T = i18n.createComponent();

// Firefox does not support named group yet
// const SendPath = new RegExp('/account/(?<address>\\w+)/(?<action>send)')
// const DelegatePath = new RegExp('/validators?/(?<address>\\w+)/(?<action>delegate)')
// const WithdrawPath = new RegExp('/account/(?<action>withdraw)')

const SendPath = new RegExp('/account/(\\w+)/(send)')
const DelegatePath = new RegExp('/validators?/(\\w+)/(delegate)')
const WithdrawPath = new RegExp('/account/(withdraw)')

const getUser = () => localStorage.getItem(CURRENTUSERADDR)

export default class Header extends Component {
    constructor(props) {
        super(props);

        this.toggle = this.toggle.bind(this);

        this.state = {
            isOpen: false,
            networks: "",
            version: "-"
        };
    }

    toggle() {
        this.setState({
            isOpen: !this.state.isOpen
        }, ()=>{
            // console.log(this.state.isOpen);
        });
    }

    toggleSignIn = (value) => {
        this.setState(( prevState) => {
            return {isSignInOpen: value!=undefined?value:!prevState.isSignInOpen}
        })
    }

    handleLanguageSwitch(lang, e) {
        i18n.setLocale(lang)
    }

    componentDidMount(){
        const url = Meteor.settings.public.networks
        if (url){
            try{
                HTTP.get(url, null, (error, result) => {
                    if (result.statusCode == 200){
                        let networks = JSON.parse(result.content);
                        if (networks.length > 0){
                            this.setState({
                                networks: <DropdownMenu>{
                                    networks.map((network, i) => {
                                        return <span key={i}>
                                            <DropdownItem header><img src={network.logo} /> {network.name}</DropdownItem>
                                            {network.links.map((link, k) => {
                                                return <DropdownItem key={k} disabled={link.chain_id == Meteor.settings.public.chainId}>
                                                    <a href={link.url} target="_blank">{link.chain_id} <Badge size="xs" color="secondary">{link.name}</Badge></a>
                                                </DropdownItem>})}
                                            {(i < networks.length - 1)?<DropdownItem divider />:''}
                                        </span>

                                    })
                                }</DropdownMenu>
                            })
                        }
                    }
                })
            }
            catch(e){
                console.warn(e);
            }
        }

        Meteor.call('getVersion', (error, result) => {
            if (result) {
                this.setState({
                    version:result
                })
            }
        })
    }

    signOut = () => {
        localStorage.removeItem(CURRENTUSERADDR);
        localStorage.removeItem(CURRENTUSERPUBKEY);
        localStorage.removeItem(BLELEDGERCONNECTION);
        localStorage.removeItem(ADDRESSINDEX);
        this.props.refreshApp();
    }

    shouldLogin = () => {
        let pathname = this.props.location.pathname
        let groups;
        let match = pathname.match(SendPath) || pathname.match(DelegatePath)|| pathname.match(WithdrawPath);
        if (match) {
            if (match[0] === '/account/withdraw') {
                groups = {action: 'withdraw'}
            } else {
                groups = {address: match[1], action: match[2]}
            }
        }
        let params = qs.parse(this.props.location.search.substr(1))
        return groups || params.signin != undefined
    }

    handleLoginConfirmed = (success) => {
        let groups = this.shouldLogin()
        if (!groups) return
        let redirectUrl;
        let params;
        if (groups) {
            let { action, address } = groups;
            params = {action}
            switch (groups.action) {
            case 'send':
                params.transferTarget = address
                redirectUrl = `/account/${address}`
                break
            case 'withdraw':
                redirectUrl = `/account/${getUser()}`
                break;
            case 'delegate':
                redirectUrl = `/validators/${address}`
                break;
            }
        } else {
            let location = this.props.location;
            params = qs.parse(location.search.substr(1))
            redirectUrl = params.redirect?params.redirect:location.pathname;
            delete params['redirectUrl']
            delete params['signin']
        }

        let query = success?`?${qs.stringify(params)}`:'';
        this.props.history.push(redirectUrl + query)
    }

    render() {
        let signedInAddress = getUser();
        return (
            <Navbar color="primary" dark expand="lg" fixed="top" id="header">
                <NavbarBrand tag={Link} to="/"><img src="/img/BaseledgerBoxWhite.png" className="img-fluid logo"/> <span className="d-none d-xl-inline-block"><T>navbar.siteName</T>&nbsp;</span><Badge color="secondary"></Badge> </NavbarBrand>

                <SearchBar id="header-search" history={this.props.history} />
                <NavbarToggler onClick={this.toggle} />
                <Collapse isOpen={this.state.isOpen} navbar>
                    <Nav className="ml-auto text-nowrap" navbar>
                        <NavItem>
                            <NavLink tag={Link} to="/validators"><T>navbar.validators</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/blocks"><T>navbar.blocks</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/transactions"><T>navbar.transactions</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/voting-power-distribution"><T>navbar.votingPower</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink href="https://baseledger.net" target="_blank" rel="noopener noreferrer">
                                baseledger.net
                            </NavLink>
                        </NavItem >
                    </Nav>
                </Collapse>
            </Navbar>
        );
    }
}