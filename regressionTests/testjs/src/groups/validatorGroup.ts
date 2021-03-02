import {TestSuite} from "../resultClasses";
import {ValidatorV1_0Handler} from "../handlers/v1.0/validatorV1_0Handler";

/*
    ValidatorGroup handles the tests for the endpoints of the validator group
 */
export class ValidatorGroup {
    v1_0Handler: ValidatorV1_0Handler;

    constructor(v1_0Handler: ValidatorV1_0Handler) {
        this.v1_0Handler = v1_0Handler;
    }

    async handleValidatorStatistics(): Promise<Array<TestSuite>> {
        let v1_0result = await this.v1_0Handler.handleValidatorStatistics();
        return new Array<TestSuite>(v1_0result);
    }
}
