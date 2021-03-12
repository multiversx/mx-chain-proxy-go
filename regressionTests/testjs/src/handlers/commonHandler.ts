import {HttpRequestHandler} from "../httpRequestHandler";
import {IProvider} from "@elrondnetwork/erdjs/out/interface";
import {Account, Transaction, UserSecretKey} from "@elrondnetwork/erdjs";
import {Check, CheckResult, TestPhase, TestStatus} from "../resultClasses";
import {Signature} from "@elrondnetwork/erdjs/out/signature";

/*
    CommonHandler class holds the needed configuration for endpoints testing and useful methods
 */
export class CommonHandler {
    proxyURL: string;
    httpRequestHandler: HttpRequestHandler;
    privateKey: UserSecretKey;
    account: Account;
    transaction: Transaction;
    provider: IProvider;

    constructor(proxyURL: string, httpRequestHandler: HttpRequestHandler, privateKey: UserSecretKey, account: Account, transaction: Transaction, provider: IProvider) {
        this.proxyURL = proxyURL;
        this.httpRequestHandler = httpRequestHandler;
        this.privateKey = privateKey;
        this.account = account;
        this.transaction = transaction;
        this.provider = provider;
    }

    getTransactionClone(): Transaction {
        let newTx = new Transaction(this.transaction);
        newTx.signature = new Signature("");

        return newTx
    }

    runBasicTestPhaseOk(response: any, code: number): TestPhase {
        let testPhase = new TestPhase("Call API endpoint");
        testPhase.checks.push(this.runCheckHttpCode(response, code));
        testPhase.checks.push(this.runCheckEmptyErrorMessage(response));

        return testPhase
    }

    runCheckHttpCode(response: any, code: number): Check {
        let check = new Check("check that response is " + code)
        if (response.status === code) {
            check.result = new CheckResult(TestStatus.SUCCESSFUL, "expected http code");
        } else {
            check.result = new CheckResult(TestStatus.UNSUCCESSFUL, "http code is " + response.status);
        }

        return check;
    }

    runCheckEmptyErrorMessage(response: any): Check {
        let check = new Check("check that the error message is empty")

        let error = response.data.error;
        if (error == "") {
            check.result = new CheckResult(TestStatus.SUCCESSFUL, "error message is empty");
        } else {
            check.result = new CheckResult(TestStatus.UNSUCCESSFUL, "error message is " + error);
        }

        return check;
    }
}
