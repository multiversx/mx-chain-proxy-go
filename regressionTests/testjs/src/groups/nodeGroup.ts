import {TestSuite} from "../resultClasses";
import {NodeV1_0Handler} from "../handlers/v1.0/nodeV1_0Handler";

/*
    NodeGroup handles the tests for the endpoints of the node group
 */
export class NodeGroup {
    v1_0Handler: NodeV1_0Handler;

    constructor(v1_0Handler: NodeV1_0Handler) {
        this.v1_0Handler = v1_0Handler;
    }

    async handleHeartbeat(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleHeartbeat();
        return new Array<TestSuite>(v1_0result);
    }
}
