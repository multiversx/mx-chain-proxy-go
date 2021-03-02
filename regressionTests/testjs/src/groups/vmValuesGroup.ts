import {TestSuite} from "../resultClasses";
import {VmValuesV1_0Handler} from "../handlers/v1.0/vmValuesV1_0Handler";

/*
    VmValuesGroup handles the tests for the endpoints of the vm values group
 */
export class VmValuesGroup {
    v1_0Handler: VmValuesV1_0Handler;

    constructor(v1_0Handler: VmValuesV1_0Handler) {
        this.v1_0Handler = v1_0Handler;
    }

    async handleVmValuesQuery(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleVmValuesQuery();
        return new Array<TestSuite>(v1_0result);
    }

    async handleVmValuesInt(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleVmValuesInt();
        return new Array<TestSuite>(v1_0result);
    }

    async handleVmValuesString(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleVmValuesString();
        return new Array<TestSuite>(v1_0result);
    }
}
