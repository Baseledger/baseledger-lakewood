/* eslint-disable camelcase */

import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Alert, UncontrolledPopover, PopoverHeader, PopoverBody } from 'reactstrap';
import { TxIcon } from '../components/Icons.jsx';
import Activities from '../components/Activities.jsx';
import CosmosErrors from '../components/CosmosErrors.jsx';
import TimeAgo from '../components/TimeAgo.jsx';
import numbro from 'numbro';
import Coin from '/both/utils/coins.js'
import SentryBoundary from '../components/SentryBoundary.jsx';
import { Markdown } from 'react-showdown';
import i18n from 'meteor/universe:i18n';

const T = i18n.createComponent();

let showdown  = require('showdown');
showdown.setFlavor('github');

const breakWords = () => {
    setTimeout(() => {
        // Set displays to block instead of inline-block to prevent side scroll
        Array.from(document.getElementsByClassName('variable-value')).forEach(item => {
            item.style.display = 'block'
            item.firstChild.style.display = 'block'
        })

        // For some reason 'display: block' is not enough on /transactions
        Array.from(document.getElementsByClassName('string-value')).forEach(item => {
            item.style.overflowWrap = 'break-word'
        })
    }, 1000); // TODO: Get to know when all span elements/activities are loaded
}

export const TransactionRow = (props) => {
    let tx = props.tx;
    let homepage = window?.location?.pathname === '/' ? true : false;

    breakWords()

    return <SentryBoundary>
        <ul>
            <li className="transaction__info">
                {/* BLOCK HEIGHT/NUMBER */}
                <div className="transaction__item">
                    <div>
                        <i className="fas fa-database"></i>

                        <span className="">
                            <T>common.height</T>
                        </span>
                    </div>

                    {(!props.blockList)
                        ? <Link to={"/blocks/"+tx.height}>{numbro(tx.height).format("0,0")}</Link>
                        : ''
                    }
                </div>

                {/* TRANSACTION VALID */}
                <div className="transaction__item">
                    <div>
                        <i className="material-icons">check_circle</i>

                        <span className="">
                            <T>transactions.valid</T>
                        </span>
                    </div>

                    {(!tx.tx_response.code)
                        ? <TxIcon valid />
                        : <TxIcon />
                    }
                </div>

                {/* TRANSACTION HASH */}
                <div className="transaction__item">
                    <div>
                        <i className="fas fa-hashtag"></i>

                        <span>
                            <T>transactions.txHash</T>
                        </span>
                    </div>

                    <Link className="transaction__hash" to={"/transactions/"+tx.txhash}>{tx.txhash}</Link>
                </div>

                {/* TRANSACTION FEE */}
                <div className="transaction__item">
                    <div>
                        <i className="material-icons">monetization_on</i>

                        <span className="">
                            <T>transactions.fee</T>
                        </span>
                    </div>

                    {(tx?.tx?.auth_info?.fee?.amount.length > 0)
                        ? tx?.tx?.auth_info?.fee?.amount.map((fee, i) => {
                            return <span className="text-nowrap transaction__fee" key={i}>
                                {(new Coin(parseFloat(fee.amount), (fee)
                                    ? fee.denom
                                    : null
                                )).toString(6)}
                            </span>
                        })
                        : <span className="transaction__fee">1token</span>
                    }
                </div>

                {/* SCHEDULED AT */}
                <div className="transaction__item">
                    <i className="material-icons">schedule</i>

                    <span>{tx.block() ? <TimeAgo time={tx.block().time} /> : ''}</span>

                    {(tx?.tx?.body?.memo && tx?.tx?.body?.memo != "")
                        ?
                        <span>
                            <i className="material-icons ml-2 memo-button" id={"memo-"+tx.txhash}>message</i>

                            <UncontrolledPopover trigger="legacy" placement="top-start" target={"memo-"+tx.txhash}>
                                <PopoverBody><Markdown markup={tx.tx.body.memo} /></PopoverBody>
                            </UncontrolledPopover>
                        </span>
                        : ''
                    }
                </div>

                {(tx.tx_response.code)
                    ?
                    <div className="error">
                        <Alert color="danger">
                            <CosmosErrors
                                code={tx.tx_response.code}
                                codespace={tx.codespace}
                                log={tx.raw_log}
                            />
                        </Alert>
                    </div>
                    : ''
                }
            </li>

            {/* TRANSACTION LOG */}
            <li
                className={
                    `${
                        (tx.tx_response.code)
                            ? "transaction__log transaction__log--invalid"
                            : "transaction__log"
                    }
                    ${homepage ? 'transaction__log--home' : ''}`
                }
            >
                {(tx?.tx?.body?.messages && tx?.tx?.body?.messages.length > 0) ? tx?.tx?.body?.messages.map((msg,i) => {
                    return <div key={i}>
                        <Activities
                            msg={msg}
                            invalid={(!!tx.tx_response.code)}
                            events={(tx.tx_response.logs&&tx.tx_response.logs[i]) ? tx.tx_response.logs[i].events : null}
                        />
                    </div>
                }) : ''}
            </li>
        </ul>
    </SentryBoundary>
}
