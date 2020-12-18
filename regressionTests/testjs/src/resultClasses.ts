/*
    Check defines a single-point step in checking an endpoint. It holds a description and a result
 */
export class Check {
    description: string;
    result: any;

    constructor(description: string) {
        this.description = description;
    }
}

/*
    Result holds data about execution of a check
 */
export class CheckResult {
    status: TestStatus;
    message: string;

    constructor(status: TestStatus, message: string) {
        this.status = status;
        this.message = message;
    }
}

/*
    TestPhase holds a description and multiple checks
 */
export class TestPhase {
    description: string;
    checks: Check[];

    constructor(description: string) {
        this.description = description;
        this.checks = new Array<Check>();
    }
}

/*
    TestSuite holds an array of multiple test phases. Useful when testing for multiple versions
 */
export class TestSuite {
    name: string;
    testPhases: TestPhase[];
    status: TestStatus;
    response: any;
    exception: any;

    constructor(name: string, testPhases: Array<TestPhase>, status: TestStatus, response: any) {
        this.name = name;
        this.testPhases = testPhases;
        this.status = status;
        this.response = response;
        this.exception = null;
    }

    static withException(name: string, err: any): TestSuite {
        let testSuite = new this(
            name,
            new Array<TestPhase>(),
            TestStatus.UNSUCCESSFUL,
            null
        );
        testSuite.exception = err;

        return testSuite
    }

    isWithException(): boolean {
        return this.exception != null;
    }

    computeStatus() {
        this.status = TestStatus.SUCCESSFUL;
        for (let phase of this.testPhases) {
            for (let check of phase.checks) {
                if (!check.result) {
                    this.status = TestStatus.UNSUCCESSFUL;
                    return;
                }
            }
        }
    }
}

/*
    TestResult holds the suites for a given test for a given endpoint
 */
export class TestResult {
    testSuites: TestSuite[];
    status: TestStatus;

    constructor() {
        this.testSuites = new Array<TestSuite>();
        this.status = TestStatus.UNSUCCESSFUL;
    }

    computeStatus() {
        this.status = TestStatus.SUCCESSFUL;
        for (let suite of this.testSuites) {
            if (suite.status == TestStatus.UNSUCCESSFUL) {
                this.status = TestStatus.UNSUCCESSFUL;
                return
            }
        }
    }
}

/*
    TestStatus is an enum that holds the statuses for a test
 */
export enum TestStatus {
    SUCCESSFUL,
    UNSUCCESSFUL
}
