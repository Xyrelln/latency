// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {app} from '../models';
import {core} from '../models';
import {adb} from '../models';
import {fs} from '../models';
import {latencywin} from '../models';

export function CalculateLatencyByCurrentImage(arg1:number):Promise<app.WinOpLatencyResult>;

export function CalculateLatencyByImageDiff(arg1:core.ImageRectInfo,arg2:number):Promise<app.WinOpLatencyResult>;

export function CheckUpdate():Promise<app.UpdateInfo>;

export function ClearCacheData():Promise<void>;

export function DoUpdate(arg1:string):Promise<void>;

export function GetDisplay(arg1:string):Promise<adb.Display>;

export function GetFirstImageInfo():Promise<core.ImageInfo>;

export function GetImage(arg1:number):Promise<app.GetImageResp>;

export function GetImageFiles():Promise<Array<string>>;

export function GetPhysicalSize(arg1:string):Promise<adb.Display>;

export function GetVersionInfo():Promise<app.VersionInfo>;

export function InputSwipe(arg1:string,arg2:adb.SwipeEvent):Promise<Error>;

export function IsAppReady():Promise<Error>;

export function IsPointerLocationOn(arg1:string):Promise<boolean>;

export function ListCaptureWindows():Promise<Array<string>>;

export function ListDevices():Promise<Array<adb.Device>>;

export function ListRecords():Promise<Array<fs.RecordFile>>;

export function OpenImageInExplorer(arg1:number):Promise<void>;

export function SetAutoSwipeOff():Promise<Error>;

export function SetAutoSwipeOn(arg1:adb.SwipeEvent,arg2:number):Promise<Error>;

export function SetPointerLocationOff(arg1:string):Promise<Error>;

export function SetPointerLocationOn(arg1:string):Promise<Error>;

export function Start(arg1:string,arg2:number):Promise<Error>;

export function StartAnalyse(arg1:core.ImageRectInfo,arg2:number):Promise<Error>;

export function StartRecord(arg1:string):Promise<Error>;

export function StartTransform():Promise<Error>;

export function StartWinOpLatency(arg1:latencywin.Config):Promise<number>;

export function StartWithVideo(arg1:string):Promise<Error>;

export function StopRunner():Promise<Error>;

export function StopScrcpyServer(arg1:string):Promise<Error>;

export function StopTransform():Promise<Error>;

export function Transform(arg1:string):Promise<Error>;

export function UploadFile(arg1:string):Promise<Error>;
