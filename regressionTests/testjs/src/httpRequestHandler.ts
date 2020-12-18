import {Err} from "@elrondnetwork/erdjs/out";

const axios = require('axios');

/*
    HttpRequestHandler is the class that is responsible for handling HTTP requests
 */
export class HttpRequestHandler {
    private proxyURL: string = "";

    constructor(url: string) {
        this.proxyURL = url;
    }

    async doGetRequest(address: string): Promise<any> {
        try {
            return await axios.get(address);
        } catch (error) {
            if (error.response.status == 500 || error.response.status == 400) {
                return error.response;
            }

            console.error(error);
            throw new Err(error)
        }
    }

    async doPostRequest(address: string, payload: string): Promise<any> {
        try {
            return await axios.post(address, payload);
        } catch (error) {
            if (error.response.status == 500 || error.response.status == 400) {
                return error.response;
            }
            throw new Err(error)
        }
    }
}
