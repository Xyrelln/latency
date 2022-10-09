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

export namespace adb {
	
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

}

export namespace core {
	
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

}

