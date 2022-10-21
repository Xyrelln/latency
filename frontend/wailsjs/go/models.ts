export namespace fs {
	
	export class RecordFile {
	    dir_name: string;
	    file_path: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new RecordFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dir_name = source["dir_name"];
	        this.file_path = source["file_path"];
	        this.size = source["size"];
	    }
	}

}

export namespace latencywin {
	
	export class InputConf {
	    type: string;
	    isAuto: boolean;
	    keyTap?: string;
	
	    static createFrom(source: any = {}) {
	        return new InputConf(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.isAuto = source["isAuto"];
	        this.keyTap = source["keyTap"];
	    }
	}
	export class Config {
	    inputCconf?: InputConf;
	    imageDiff_threshold: number;
	    frames?: number;
	    startKey: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputCconf = this.convertValues(source["inputCconf"], InputConf);
	        this.imageDiff_threshold = source["imageDiff_threshold"];
	        this.frames = source["frames"];
	        this.startKey = source["startKey"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace app {
	
	export class WinOpLatencyResult {
	    latency: number;
	    responseIndex: number;
	    responseTime: number;
	
	    static createFrom(source: any = {}) {
	        return new WinOpLatencyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.latency = source["latency"];
	        this.responseIndex = source["responseIndex"];
	        this.responseTime = source["responseTime"];
	    }
	}
	export class UpdateInfo {
	    latestVersion: string;
	    needUpdate: boolean;
	    err?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.latestVersion = source["latestVersion"];
	        this.needUpdate = source["needUpdate"];
	        this.err = source["err"];
	    }
	}
	export class GetImageResp {
	    length: number;
	    currentIndex: number;
	    screenshotTime: string;
	    imageFilePath: string;
	    imageWidth: number;
	    imageHeight: number;
	
	    static createFrom(source: any = {}) {
	        return new GetImageResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.length = source["length"];
	        this.currentIndex = source["currentIndex"];
	        this.screenshotTime = source["screenshotTime"];
	        this.imageFilePath = source["imageFilePath"];
	        this.imageWidth = source["imageWidth"];
	        this.imageHeight = source["imageHeight"];
	    }
	}
	export class VersionInfo {
	    version: string;
	    commitShortSHA: string;
	    buildTimestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new VersionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.commitShortSHA = source["commitShortSHA"];
	        this.buildTimestamp = source["buildTimestamp"];
	    }
	}

}

export namespace core {
	
	export class ImageRectInfo {
	    x: number;
	    y: number;
	    w: number;
	    h: number;
	    preview_width: number;
	    preview_height: number;
	    source_width: number;
	    source_height: number;
	
	    static createFrom(source: any = {}) {
	        return new ImageRectInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.w = source["w"];
	        this.h = source["h"];
	        this.preview_width = source["preview_width"];
	        this.preview_height = source["preview_height"];
	        this.source_width = source["source_width"];
	        this.source_height = source["source_height"];
	    }
	}
	export class ImageInfo {
	    path: string;
	    width: number;
	    height: number;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ImageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.size = source["size"];
	    }
	}

}

export namespace adb {
	
	export class SwipeEvent {
	    sx: number;
	    sy: number;
	    dx: number;
	    dy: number;
	    speed: number;
	
	    static createFrom(source: any = {}) {
	        return new SwipeEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sx = source["sx"];
	        this.sy = source["sy"];
	        this.dx = source["dx"];
	        this.dy = source["dy"];
	        this.speed = source["speed"];
	    }
	}
	export class Device {
	    Serial: string;
	    State: number;
	    abi: string;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Serial = source["Serial"];
	        this.State = source["State"];
	        this.abi = source["abi"];
	    }
	}
	export class Display {
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new Display(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}

}

