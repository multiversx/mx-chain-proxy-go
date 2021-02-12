import {TestPhase, TestStatus, TestSuite} from "../../resultClasses";
import {CommonHandler} from "../commonHandler";

/*
    ActionsV1_0Handler is the class that will handle the node API calls for the API v1.0
 */
export class ActionsV1_0Handler {
    commonHandler: CommonHandler;

    constructor(commonHandler: CommonHandler) {
        this.commonHandler = commonHandler;
    }

    async handleReloadObserversShouldNotBeAuthorized(): Promise<TestSuite> {
        let testPhases = new Array<TestPhase>();

        let url = this.commonHandler.proxyURL + "/actions/reload-observers";
        try {
            let response = await this.commonHandler.httpRequestHandler.doPostRequest(url, "")
            let testPhase = new TestPhase("check that return code is unauthorized");
            testPhase.checks.push(this.commonHandler.runCheckHttpCode(response, 401))

            testPhases.push(testPhase)

            return new TestSuite("v1.0", testPhases, TestStatus.SUCCESSFUL, response)

        } catch (error) {
            return TestSuite.withException("v1.0", error);
        }
    }
}
