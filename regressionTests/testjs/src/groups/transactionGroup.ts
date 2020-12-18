import {TestSuite} from "../resultClasses";
import {TransactionV1_0Handler} from "../handlers/v1.0/transactionV1_0Handler";

/*
    TransactionGroup handles the tests for the endpoints of the transaction group
 */
export class TransactionGroup {
    v1_0Handler: TransactionV1_0Handler;

    constructor(v1_0Handler: TransactionV1_0Handler) {
        this.v1_0Handler = v1_0Handler;
    }

    async handleSendTransaction() {
        let v1_0result = await this.v1_0Handler.handleSendTransaction();
        return new Array<TestSuite>(v1_0result);
    }

    async handleSendMultipleTransactions(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleSendMultipleTransactions();
        return new Array<TestSuite>(v1_0result);
    }
}
