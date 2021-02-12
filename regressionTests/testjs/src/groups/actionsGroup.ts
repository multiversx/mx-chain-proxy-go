import {TestSuite} from "../resultClasses";
import {ActionsV1_0Handler} from "../handlers/v1.0/actionsV1_0Handler";

/*
    ActionsGroup handles the tests for the endpoints of the actions group
 */
export class ActionsGroup {
    v1_0Handler: ActionsV1_0Handler;

    constructor(v1_0Handler: ActionsV1_0Handler) {
        this.v1_0Handler = v1_0Handler;
    }

    async handleReloadObserversShouldNotBeAuthorized(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleReloadObserversShouldNotBeAuthorized();
        return new Array<TestSuite>(v1_0result);
    }
}
