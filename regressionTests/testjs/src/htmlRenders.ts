import {CheckResult, TestStatus, TestSuite} from "./resultClasses";

declare var $: any;

export function displayTestPhases(container: any, testSuite: TestSuite): string {
    if (testSuite.isWithException()) {
        return `<div class="card text-white bg-danger mb-3" style="max-width: 18rem;">
  <div class="card-header">Exception</div>
  <div class="card-body">
    <p class="card-text">${testSuite.exception}</p>
  </div>
</div>`
    }

    if (testSuite.response == null) {
        testSuite.response = {config: {url: ""}, status: 0, data: {"data": "", "error": ""}}
    }
    let mainContent = `
    <h5><span class="badge badge-primary">${testSuite.response.config.url} <span class="badge badge-light">${testSuite.response.status}</span></span></h5>
    <div class="row">
  <div class="col-sm-4">
    <div class="card">
      <div class="card-body">
        <h5 class="card-title">Request</h5>
        <p class="card-text"><pre>${testSuite.response.config.data}</pre></p>
      </div>
    </div>
  </div>
  <div class="col-sm-8 overflow-auto" style="max-height: 200px;">
    <div class="card">
      <div class="card-body">
        <h5 class="card-title">Response</h5>
        <p class="card-text"><pre>${JSON.stringify(testSuite.response.data, null, 2)}</pre></p>
      </div>
    </div>
  </div>
</div>
    `
    testSuite.testPhases.forEach(function (testPhase) {
        mainContent += `<h5>${testPhase.description}</h5>`;
        testPhase.checks.forEach(function (check) {
            mainContent += `
                <ul>
                    <li class="list-group-item">Description: ${check.description}</li>
                    ${displayResult(check.result)}
                </ul>`
        })
    })
    $(`#${container}`).html(mainContent);

    return mainContent;
}

export function displayResult(result: CheckResult): string {
    if (result.status == TestStatus.SUCCESSFUL) {
        return `<li class="list-group-item list-group-item-success">${result.message}</li>`
    }

    return `<li class="list-group-item list-group-item-danger">${result.message}</li>`
}

export function displayErrResponse(container: string, url: string, obj: any) {
    let alert = `
        <ul>
            <li class="list-group-item">URL: ${url}</li>
            <li class="list-group-item list-group-item-danger">${obj}</li>
        </ul>`;
    $(`#${container}`).html(alert);
}

export function displayObject(container: string, obj: any) {
    let json = JSON.stringify(obj, null, 4);
    $(`#${container}`).html(json);
}

export function displayOkResponse(container: string, url: string, obj: any) {
    let json = JSON.stringify(obj, null, 4);
    let okAlert = `
        <ul>
            <li class="list-group-item">URL: ${url}</li>
            <li class="list-group-item list-group-item-success">${json}</li>
        </ul>`;
    $(`#${container}`).html(okAlert);
}

export function displayTestsSuites(container: any, suites: Array<TestSuite>) {
    let content = ""
    suites.forEach(function (value) {
        content += `<h3>${value.name}</h3>`;
        content += displayTestPhases(container, value);
    })

    $(`#${container}`).html(content);
}
