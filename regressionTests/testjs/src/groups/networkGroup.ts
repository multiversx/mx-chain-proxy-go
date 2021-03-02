import {TestSuite} from "../resultClasses";
import {NetworkV1_0Handler} from "../handlers/v1.0/networkV1_0Handler";

/*
    NetworkGroup handles the tests for the endpoints of the network group
 */
export class NetworkGroup {
    v1_0Handler: NetworkV1_0Handler;

    constructor(v1_0Handler: NetworkV1_0Handler) {
        this.v1_0Handler = v1_0Handler;
    }

    async handleNetworkConfig(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleNetworkConfig();
        return new Array<TestSuite>(v1_0result);
    }

    async handleNetworkStatus(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleNetworkStatus();
        return new Array<TestSuite>(v1_0result);
    }

    async handlerNetworkTotalStaked(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleNetworkTotalStaked();
        return new Array<TestSuite>(v1_0result);
    }
}
