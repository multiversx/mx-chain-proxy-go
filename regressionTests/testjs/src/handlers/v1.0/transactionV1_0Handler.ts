import {Nonce, SimpleSigner} from "@elrondnetwork/erdjs/out";
import {Check, CheckResult, TestPhase, TestStatus, TestSuite} from "../../resultClasses";
import {CommonHandler} from "../commonHandler";

/*
    TransactionV1_0Handler is the class that will handle the transaction API calls for the API v1.0
 */
export class TransactionV1_0Handler {
    commonHandler: CommonHandler;

    constructor(commonHandler: CommonHandler) {
        this.commonHandler = commonHandler;
    }

    async handleSendTransaction(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let testPhase = new TestPhase("Phase 0: send a transaction");
        try {
            let signer = new SimpleSigner(this.commonHandler.privateKey);
            await signer.sign(this.commonHandler.transaction);
            let check = new Check("sign and marshal a transaction");
            let txJson = JSON.stringify(this.commonHandler.transaction.toPlainObject(), null, 4);
            check.result = new CheckResult(TestStatus.SUCCESSFUL, txJson);
            testPhase.checks.push(check);
            testPhases.push(testPhase);

            let url = this.commonHandler.proxyURL + "/transaction/send"
            try {
                let response = await this.commonHandler.httpRequestHandler.doPostRequest(url, txJson)

                testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

                return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
            } catch (error) {
                return TestSuite.withException("v1.0", error);
            }

        } catch (error) {
            return TestSuite.withException("v1.0", error);
        }
    }

    async handleSendMultipleTransactions(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let testPhase = new TestPhase("Phase 0: send a transaction");
        try {
            let signer = new SimpleSigner(this.commonHandler.privateKey);
            await signer.sign(this.commonHandler.transaction);
            let check = new Check("sign and marshal transactions");
            let txJson1 = JSON.stringify(this.commonHandler.transaction.toPlainObject(), null, 4);
            let tx2 = this.commonHandler.transaction;
            tx2.setNonce(new Nonce(this.commonHandler.transaction.nonce.valueOf() + 1));
            await signer.sign(tx2);
            let txJson2 = JSON.stringify(tx2.toPlainObject(), null, 4);

            let txsJson = `[${txJson1}, ${txJson2}]`;
            console.log(txsJson);
            check.result = new CheckResult(TestStatus.SUCCESSFUL, txsJson);
            testPhase.checks.push(check);
            testPhases.push(testPhase);

            let url = this.commonHandler.proxyURL + "/transaction/send-multiple"
            try {
                let response = await this.commonHandler.httpRequestHandler.doPostRequest(url, txsJson)

                testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

                return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
            } catch (error) {
                return new TestSuite("v1.0", testPhases, TestStatus.UNSUCCESSFUL, null)
            }

        } catch (error) {
            return new TestSuite("v1.0", testPhases, TestStatus.UNSUCCESSFUL, null)
        }
    }
}
