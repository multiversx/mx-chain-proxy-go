import {Err} from "@elrondnetwork/erdjs/out";
import {TestPhase, TestStatus, TestSuite} from "../../resultClasses";
import {displayErrResponse} from "../../htmlRenders";
import {CommonHandler} from "../commonHandler";

/*
    VmValuesV1_0Handler is the class that will handle the vm values API calls for the API v1.0
 */
export class VmValuesV1_0Handler {
    commonHandler: CommonHandler;

    constructor(commonHandler: CommonHandler) {
        this.commonHandler = commonHandler;
    }

    async handleVmValuesQuery(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/vm-values/query";
        try {
            let body = `{"scAddress":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7","funcName": "getQueueSize","args": []}`
            let response = await this.commonHandler.httpRequestHandler.doPostRequest(url, body)

            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))
            console.log(testPhases);

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            displayErrResponse("LoadVmValuesQueryOutput", url, Err.html(error))
            return new TestSuite("v1.0", testPhases, TestStatus.UNSUCCESSFUL, null)
        }
    }

    async handleVmValuesInt(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/vm-values/int";
        try {
            let body = `{"scAddress":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7","funcName": "getQueueSize","args": []}`
            let response = await this.commonHandler.httpRequestHandler.doPostRequest(url, body)

            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            displayErrResponse("LoadVmValuesIntOutput", url, Err.html(error))
            return new TestSuite("v1.0", testPhases, TestStatus.UNSUCCESSFUL, null)
        }
    }

    async handleVmValuesString(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/vm-values/string";
        try {
            let body = `{"scAddress":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7","funcName": "getQueueSize","args": []}`
            let response = await this.commonHandler.httpRequestHandler.doPostRequest(url, body)

            testPhases.push(this.commonHandler.runBasicTestPhaseOk(response, 200))

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)
        } catch (error) {
            displayErrResponse("LoadVmValuesStringOutput", url, Err.html(error))
            return new TestSuite("v1.0", testPhases, TestStatus.UNSUCCESSFUL, null)
        }
    }
}
