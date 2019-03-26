
import * as React from 'react';
export class RemoteClient extends React.Component{

    private url ="";

	constructor(props:any) {
		super(props)
		this.url = "/api/"
	}

    public sign(payload: {name:string, date:string}) {

        const url = this.url + "sign"
        const body = payload
        
        const axios = require('axios');
		const response = axios.create({
		    headers:this.headers(),
		    responseType: 'json',
			timeout: 15000

			}).post(url,body).then((r)=>r.data).then(data=>{
				console.log(data)
				return {signed: data.signed, unsigned: data.unsigned}
			}).catch(e=>{
				// console.log("OOpps...",e)
				if (e.response) {
					return {status:"error",reason:e.response.data}
				}
				return {status:"error",reason:e.message}
			})
        return response
    }

	public get() {

		const url = this.url + "names"
        
        const axios = require('axios');
        const response = axios.create({
                headers: this.headers(),
				responseType: 'json',
				timeout: 15000

			}).get(url).then(x=>x.data).then(data=>{
				for (const i in data.unsigned) {
					console.log("inside", data.unsigned[i])
				}

				return {signed: data.signed, unsigned: data.unsigned}

			}).catch(e=>{
				if (e.response) {
					return {status:"error",reason:e.response.data}
				}
				return {status:"error",reason:e.message}
            })

		return response

	}

	private headers(opt: { session?:string, timeout?: number} = {}) {

		const headers = {
			"Content-type": "application/json",
		} as any;

		return headers;
	}


};

