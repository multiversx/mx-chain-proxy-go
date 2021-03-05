import {
    Account,
    Err,
    NetworkConfig,
    ProxyProvider,
    Transaction,
    UserSecretKey,
    UserSigner,
} from "@elrondnetwork/erdjs";
import {HttpRequestHandler} from "./httpRequestHandler"
import {ValidatorGroup} from "./groups/validatorGroup";
import {AddressGroup,} from "./groups/addressGroup";
import {NodeGroup} from "./groups/nodeGroup";
import {TransactionGroup} from "./groups/transactionGroup";
import {displayTestsSuites} from "./htmlRenders";
import {NetworkV1_0Handler} from "./handlers/v1.0/networkV1_0Handler";
import {NetworkGroup} from "./groups/networkGroup";
import {AddressV1_0Handler} from "./handlers/v1.0/addressV1_0Handler";
import {CommonHandler} from "./handlers/commonHandler";
import {ValidatorV1_0Handler} from "./handlers/v1.0/validatorV1_0Handler";
import {NodeV1_0Handler} from "./handlers/v1.0/nodeV1_0Handler";
import {VmValuesGroup} from "./groups/vmValuesGroup";
import {VmValuesV1_0Handler} from "./handlers/v1.0/vmValuesV1_0Handler";
import {TransactionV1_0Handler} from "./handlers/v1.0/transactionV1_0Handler";
import {ActionsV1_0Handler} from "./handlers/v1.0/actionsV1_0Handler";
import {ActionsGroup} from "./groups/actionsGroup";

declare var $: any;
let commonHandler: CommonHandler;
let addressGroup: AddressGroup;
let actionsGroup: ActionsGroup;
let nodeGroup: NodeGroup;
let networkGroup: NetworkGroup;
let transactionGroup: TransactionGroup;
let validatorGroup: ValidatorGroup;
let vmValuesGroup: VmValuesGroup;

$(async function () {
    $("#LoadAccountButton").click(async function () {
        try {
            let provider = new ProxyProvider(getProxyUrl());
            let transaction = new Transaction();
            let httpRequestHandler = new HttpRequestHandler(getProxyUrl());

            try {
                await NetworkConfig.getDefault().sync(provider);
            } catch (error) {
                onError(error);
            }

            let signer = new UserSigner(getPrivateKey());
            let account = new Account(signer.getAddress());
            await account.sync(provider);

            $("#AccountAddress").text(account.address.bech32());
            $("#AccountNonce").text(account.nonce.valueOf());
            $("#AccountBalance").text(account.balance.toString());
            commonHandler = new CommonHandler(getProxyUrl(), httpRequestHandler, getPrivateKey(), account, transaction, provider)

            let addressGroupV1_0Handler = new AddressV1_0Handler(commonHandler);
            addressGroup = new AddressGroup(addressGroupV1_0Handler);

            let actionsGroupV1_0Handler = new ActionsV1_0Handler(commonHandler);
            actionsGroup = new ActionsGroup(actionsGroupV1_0Handler);

            let nodeV1_0Handler = new NodeV1_0Handler(commonHandler);
            nodeGroup = new NodeGroup(nodeV1_0Handler);

            let networkGroupV1_0Handler = new NetworkV1_0Handler(commonHandler);
            networkGroup = new NetworkGroup(networkGroupV1_0Handler);

            let transactionV1_0Handler = new TransactionV1_0Handler(commonHandler);
            transactionGroup = new TransactionGroup(transactionV1_0Handler);

            let validatorV1_0Handler = new ValidatorV1_0Handler(commonHandler);
            validatorGroup = new ValidatorGroup(validatorV1_0Handler);

            let vmValuesV1_0Handler = new VmValuesV1_0Handler(commonHandler);
            vmValuesGroup = new VmValuesGroup(vmValuesV1_0Handler);

        } catch (error) {
            onError(error);
        }
    });

    $("#LoadAccountDataShouldWork").click(async function () {
        try {
            let response = await addressGroup.handleAddressShouldWork()
            displayTestsSuites("LoadAccountDataShouldWorkOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    $("#LoadAccountDataShouldErr").click(async function () {
        try {
            let response = await addressGroup.handleAddressShouldErr();
            displayTestsSuites("LoadAccountDataShouldErrOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    $("#LoadAccountNonce").click(async function () {
        try {
            let response = await addressGroup.handleAddressNonce();
            displayTestsSuites("LoadAccountNonceOutput", response);
        } catch (error) {
            onError(error);
        }

    });

    $("#LoadAccountBalance").click(async function () {
        try {
            let response = await addressGroup.handleAddressBalance();
            displayTestsSuites("LoadAccountBalanceOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    $("#LoadAccountShard").click(async function () {
        try {
            let response = await addressGroup.handleAddressShard();
            displayTestsSuites("LoadAccountShardOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    $("#LoadAccountTransactions").click(async function () {
        try {
            let response = await addressGroup.handleAddressTransactions();
            displayTestsSuites("LoadAccountTransactionsOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    // -- Actions group

    $("#LoadReloadObservers").click(async function () {
        try {
            let response = await actionsGroup.handleReloadObserversShouldNotBeAuthorized();
            displayTestsSuites("LoadReloadObserversOutput", response);
        } catch (error) {
            onError(error);
        }
    });


    // -- Network group

    $("#LoadNetworkConfig").click(async function () {
        try {
            let response = await networkGroup.handleNetworkConfig();
            displayTestsSuites("LoadNetworkConfigOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    $("#LoadNetworkStatus").click(async function () {
        try {
            let response = await networkGroup.handleNetworkStatus();
            displayTestsSuites("LoadNetworkStatusOutput", response);
        } catch (error) {
            onError(error);
        }
    });


    // -- Node group

    $("#LoadHeartBeat").click(async function () {
        try {
            let response = await nodeGroup.handleHeartbeat();
            displayTestsSuites("LoadHeartBeatOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    // -- Transaction group

    $("#LoadSendTransaction").click(async function () {
        try {
            let response = await transactionGroup.handleSendTransaction();
            displayTestsSuites("LoadSendTransactionOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    $("#LoadSendMultipleTransactions").click(async function () {
        try {
            let response = await transactionGroup.handleSendMultipleTransactions();
            displayTestsSuites("LoadSendMultipleTransactionsOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    // -- Validator group

    $("#LoadValidatorStatistics").click(async function () {
        try {
            let response = await validatorGroup.handleValidatorStatistics();
            displayTestsSuites("LoadValidatorStatisticsOutput", response);
        } catch (error) {
            onError(error);
        }
    });


    // -- VM Values group

    $("#LoadVmValuesQuery").click(async function () {
        try {
            let response = await vmValuesGroup.handleVmValuesQuery();
            displayTestsSuites("LoadVmValuesQueryOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    $("#LoadVmValuesInt").click(async function () {
        try {
            let response = await vmValuesGroup.handleVmValuesInt();
            displayTestsSuites("LoadVmValuesIntOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    $("#LoadVmValuesString").click(async function () {
        try {
            let response = await vmValuesGroup.handleVmValuesString();
            displayTestsSuites("LoadVmValuesStringOutput", response);
        } catch (error) {
            onError(error);
        }
    });

    // - Run all tests
    $("#RunAllTests").click(async function () {
        await runAllTests();
    })
});

function onError(error: Error) {
    let html = Err.html(error);
    $("#ErrorModal .error-text").html(html);
    $("#ErrorModal").modal("show");
}

function getProxyUrl(): string {
    return $("#ProxyInput").val();
}

function getPrivateKey(): UserSecretKey {
    return UserSecretKey.fromString($("#PrivateKeyInput").val().trim());
}

async function runAllTests() {
    /*
        Tests who trigger transactions sending requests are skipped because of nonce issues as these tests use
        the same account.
     */
    $("#LoadAccountDataShouldWork").click();
    $("#LoadAccountDataShouldErr").click();
    $("#LoadAccountNonce").click();
    $("#LoadAccountBalance").click();
    $("#LoadAccountShard").click();
    //$("#LoadAccountTransactions").click();

    $("#LoadHeartBeat").click();

    $("#LoadNetworkConfig").click();
    $("#LoadNetworkStatus").click();

    //$("#LoadSendTransaction").click();
    //$("#LoadSendMultipleTransactions").click();

    $("#LoadValidatorStatistics").click();

    $("#LoadVmValuesQuery").click();
    $("#LoadVmValuesInt").click();
    $("#LoadVmValuesString").click();
    $("#LoadNetworkTotalStaked").click();
}
