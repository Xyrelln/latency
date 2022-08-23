export namespace adb {
	
	export class Device {
	    Serial: string;
	    State: number;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Serial = source["Serial"];
	        this.State = source["State"];
	    }
	}

}

