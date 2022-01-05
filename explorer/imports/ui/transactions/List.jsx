import React, { Component } from 'react';
import { Row, Col, Spinner } from 'reactstrap';
import { TransactionRow } from './TransactionRow.jsx';
import i18n from 'meteor/universe:i18n';

const T = i18n.createComponent();

export default class Transactions extends Component{
    constructor(props){
        super(props);
        this.state = {
            txs: "",
            homepage:  window?.location?.pathname === '/' ? true : false,
        }
    }

    componentDidUpdate(prevProps){
        if (this.props != prevProps){
            if (this.props.transactions.length > 0){
                this.setState({
                    txs: this.props.transactions.map((tx, i) => {
                        return <TransactionRow
                            key={i}
                            index={i}
                            tx={tx}
                        />
                    })
                })
            }
        }
    }

    getItemSize = index => (117 * (this.state.txs[index]?.props?.tx?.tx?.body?.messages.length))

    render(){
        if (this.props.loading){
            return <Spinner type="grow" color="primary" />
        } else if (!this.props.transactionsExist){
            return <div><T>transactions.notFound</T></div>
        } else {
            return <div className={`transactions-list ${this.state.homepage ? 'transactions-list--home' : ''}`}>
                <div className="activities">
                    <i className="material-icons">message</i>

                    <p className="activities__title">
                        <T>transactions.activities</T>
                    </p>
                </div>

                <div className={`wrapper ${this.state.homepage ? 'wrapper--home' : ''}`}>
                    {this.state.txs}
                </div>
            </div>
        }
    }
}
