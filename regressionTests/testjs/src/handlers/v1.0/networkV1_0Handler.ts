import {TestPhase, TestStatus, TestSuite} from "../../resultClasses";
import {CommonHandler} from "../commonHandler";

/*
    NetworkV1_0Handler is the class that will handle the network API calls for the API v1.0
 */
export class NetworkV1_0Handler {
    commonHandler: CommonHandler;

    constructor(commonHandler: CommonHandler) {
        this.commonHandler = commonHandler;
    }

    async handleNetworkConfig(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/network/config";
        try {
            let response = await this.commonHandler.httpRequestHandler.doGetRequest(url)
            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)

        } catch (error) {
            return TestSuite.withException("v1.0", error);
        }
    }

    async handleNetworkStatus(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/network/status/0";
        try {
            let response = await this.commonHandler.httpRequestHandler.doGetRequest(url)
            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            return TestSuite.withException("v1.0", error);
        }
    }
}
