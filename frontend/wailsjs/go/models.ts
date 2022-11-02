export namespace lighttestservice {
	
	export class UserInfo {
	    username: string;
	    client: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new UserInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.client = source["client"];
	        this.description = source["description"];
	    }
	}

}

export namespace adb {
	
	export class Display {
	    width: number;
	    height: number;
	    app_width: number;
	    app_height: number;
	
	    static createFrom(source: any = {}) {
	        return new Display(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	        this.app_width = source["app_width"];
	        this.app_height = source["app_height"];
	    }
	}
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
	export class TapEvent {
	    x: number;
	    y: number;
	
	    static createFrom(source: any = {}) {
	        return new TapEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	    }
	}
	export class Device {
	    serial: string;
	    state: number;
	    abi?: string;
	    usb?: string;
	    product: string;
	    model: string;
	    device: string;
	    transport_id: number;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.serial = source["serial"];
	        this.state = source["state"];
	        this.abi = source["abi"];
	        this.usb = source["usb"];
	        this.product = source["product"];
	        this.model = source["model"];
	        this.device = source["device"];
	        this.transport_id = source["transport_id"];
	    }
	}

}

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
	    keyTap: string;
	
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
	    inputConf?: InputConf;
	    captureWindow: string;
	    frames?: number;
	    startKey: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputConf = this.convertValues(source["inputConf"], InputConf);
	        this.captureWindow = source["captureWindow"];
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
	
	export class userSecret {
	    username: string;
	    key: string;
	
	    static createFrom(source: any = {}) {
	        return new userSecret(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.key = source["key"];
	    }
	}
	export class Threshold {
	    pointer_threshold: number;
	    black_white_threshold: number;
	    scene_threshold: number;
	
	    static createFrom(source: any = {}) {
	        return new Threshold(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pointer_threshold = source["pointer_threshold"];
	        this.black_white_threshold = source["black_white_threshold"];
	        this.scene_threshold = source["scene_threshold"];
	    }
	}
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
	export class UserAction {
	    auto: boolean;
	    type: string;
	    x: number;
	    y: number;
	    tx: number;
	    ty: number;
	    speed: number;
	
	    static createFrom(source: any = {}) {
	        return new UserAction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.auto = source["auto"];
	        this.type = source["type"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.tx = source["tx"];
	        this.ty = source["ty"];
	        this.speed = source["speed"];
	    }
	}
	export class CropInfo {
	    top: number;
	    left: number;
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new CropInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.top = source["top"];
	        this.left = source["left"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}
	export class DeviceInfo {
	    device_name: string;
	    screen_width: number;
	    screen_height: number;
	
	    static createFrom(source: any = {}) {
	        return new DeviceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.device_name = source["device_name"];
	        this.screen_width = source["screen_width"];
	        this.screen_height = source["screen_height"];
	    }
	}
	export class UserScene {
	    name: string;
	    key: string;
	    device: DeviceInfo;
	    crop_coordinate: CropInfo;
	    crop_touch_coordinate: CropInfo;
	    action: UserAction;
	
	    static createFrom(source: any = {}) {
	        return new UserScene(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.key = source["key"];
	        this.device = this.convertValues(source["device"], DeviceInfo);
	        this.crop_coordinate = this.convertValues(source["crop_coordinate"], CropInfo);
	        this.crop_touch_coordinate = this.convertValues(source["crop_touch_coordinate"], CropInfo);
	        this.action = this.convertValues(source["action"], UserAction);
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

