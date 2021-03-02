import {TestSuite} from "../resultClasses";
import {AddressV1_0Handler} from "../handlers/v1.0/addressV1_0Handler";

/*
    AddressGroup handles the test for the endpoints of the address group
 */
export class AddressGroup {
    v1_0Handler: AddressV1_0Handler;

    constructor(v1_0Handler: AddressV1_0Handler) {
        this.v1_0Handler = v1_0Handler;
    }

    async handleAddressShouldWork(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleAddressShouldWork();
        return new Array<TestSuite>(v1_0result);
    }

    async handleAddressShouldErr(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleAddressShouldErr();
        return new Array<TestSuite>(v1_0result);
    }

    async handleAddressNonce(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleAddressNonce();
        return new Array<TestSuite>(v1_0result);
    }

    async handleAddressBalance(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleAddressBalance();
        return new Array<TestSuite>(v1_0result);
    }

    async handleAddressShard(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleAddressShard();
        return new Array<TestSuite>(v1_0result);
    }

    async handleAddressTransactions(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleAddressTransactions();
        return new Array<TestSuite>(v1_0result);
    }
}
