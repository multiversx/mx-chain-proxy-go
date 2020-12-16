import {Check, CheckResult, TestPhase, TestStatus, TestSuite} from "../../resultClasses";
import {Account, SimpleSigner} from "@elrondnetwork/erdjs";
import {CommonHandler} from "../commonHandler";

/*
    AddressV1_0Handler is the class that will handle the address API calls for the API v1.0
 */
export class AddressV1_0Handler {
    commonHandler: CommonHandler;

    constructor(commonHandler: CommonHandler) {
        this.commonHandler = commonHandler;
    }

    async handleAddressShouldWork(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let signer = new SimpleSigner(this.commonHandler.privateKey);
        let account = new Account(signer.getAddress());
        let url = this.commonHandler.proxyURL + "/address/" + account.address.bech32();
        try {
            let response = await this.commonHandler.httpRequestHandler.doGetRequest(url);
            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))
            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            return TestSuite.withException("v1.0", error)
        }
    }

    async handleAddressShouldErr(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/address/erd1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsmsglzz"
        try {
            let response = await this.commonHandler.httpRequestHandler.doGetRequest(url)

            let testPhase = new TestPhase("Call API endpoint");
            let check = this.commonHandler.runCheckHttpCode(response, 500)
            testPhase.checks.push(check);

            let check2 = new Check("check that the error message reports a bech32 issue");
            if (response.data.error.includes("checksum failed")) {
                check2.result = new CheckResult(TestStatus.SUCCESSFUL, response.data.error)
            } else {
                check2.result = new CheckResult(TestStatus.SUCCESSFUL, response.data.error)
            }
            testPhase.checks.push(check2);

            testPhases.push(testPhase);

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            return TestSuite.withException("v1.0", error)
        }
    }

    async handleAddressNonce(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/address/" + this.commonHandler.account.address.bech32() + "/nonce";
        try {
            let response = await this.commonHandler.httpRequestHandler.doGetRequest(url)
            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            return TestSuite.withException("v1.0", error)
        }
    }

    async handleAddressBalance(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/address/" + this.commonHandler.account.address.bech32() + "/balance";
        try {
            let response = await this.commonHandler.httpRequestHandler.doGetRequest(url)
            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            return TestSuite.withException("v1.0", error)
        }
    }

    async handleAddressShard(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/address/" + this.commonHandler.account.address.bech32() + "/shard";
        try {
            let response = await this.commonHandler.httpRequestHandler.doGetRequest(url)
            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            return TestSuite.withException("v1.0", error)
        }
    }

    async handleAddressTransactions(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let testPhase = new TestPhase("Phase 0: send a transaction");
        try {
            let transaction = this.commonHandler.getTransactionClone();

            await this.commonHandler.account.sync(this.commonHandler.provider);
            transaction.setNonce(this.commonHandler.account.nonce);
            let signer = new SimpleSigner(this.commonHandler.privateKey);
            await signer.sign(transaction);
            let check = new Check("sign and marshal a transactions");
            let txJson = JSON.stringify(transaction.toPlainObject(), null, 4);
            check.result = new CheckResult(TestStatus.SUCCESSFUL, txJson);
            testPhase.checks.push(check);

            let checkSendTx = new Check("send transaction to network");
            try {
                let transactionHash = await transaction.send(this.commonHandler.provider);
                console.log(transactionHash);
                checkSendTx.result = new CheckResult(TestStatus.SUCCESSFUL, "tx hash: " + transactionHash)
            } catch (error) {
                checkSendTx.result = new CheckResult(TestStatus.UNSUCCESSFUL, error)
            }
            testPhase.checks.push(checkSendTx);

            testPhases.push(testPhase);
            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, null)
        } catch (error) {
            return TestSuite.withException("v1.0", error)
        }
    }
}
