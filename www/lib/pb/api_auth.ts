// @generated by protobuf-ts 2.9.3
// @generated from protobuf file "api_auth.proto" (package "api_auth", syntax proto3)
// tslint:disable
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { Status } from "./common";
/**
 * @generated from protobuf message api_auth.LoginRequest
 */
export interface LoginRequest {
    /**
     * @generated from protobuf field: optional string username = 1;
     */
    username?: string;
    /**
     * @generated from protobuf field: optional string password = 2;
     */
    password?: string;
}
/**
 * @generated from protobuf message api_auth.LoginResponse
 */
export interface LoginResponse {
    /**
     * @generated from protobuf field: optional common.Status status = 1;
     */
    status?: Status;
    /**
     * @generated from protobuf field: optional string token = 2;
     */
    token?: string;
}
/**
 * @generated from protobuf message api_auth.RegisterRequest
 */
export interface RegisterRequest {
    /**
     * @generated from protobuf field: optional string username = 1;
     */
    username?: string;
    /**
     * @generated from protobuf field: optional string password = 2;
     */
    password?: string;
    /**
     * @generated from protobuf field: optional string email = 3;
     */
    email?: string;
}
/**
 * @generated from protobuf message api_auth.RegisterResponse
 */
export interface RegisterResponse {
    /**
     * @generated from protobuf field: optional common.Status status = 1;
     */
    status?: Status;
}
/**
 * @generated from protobuf message api_auth.APIPermission
 */
export interface APIPermission {
    /**
     * @generated from protobuf field: optional string method = 1;
     */
    method?: string;
    /**
     * @generated from protobuf field: optional string path = 2;
     */
    path?: string;
}
/**
 * @generated from protobuf message api_auth.SignTokenRequest
 */
export interface SignTokenRequest {
    /**
     * @generated from protobuf field: optional int64 expires_in = 1;
     */
    expiresIn?: bigint;
    /**
     * @generated from protobuf field: repeated api_auth.APIPermission permissions = 2;
     */
    permissions: APIPermission[];
}
/**
 * @generated from protobuf message api_auth.SignTokenResponse
 */
export interface SignTokenResponse {
    /**
     * @generated from protobuf field: optional common.Status status = 1;
     */
    status?: Status;
    /**
     * @generated from protobuf field: optional string token = 2;
     */
    token?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class LoginRequest$Type extends MessageType<LoginRequest> {
    constructor() {
        super("api_auth.LoginRequest", [
            { no: 1, name: "username", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "password", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<LoginRequest>): LoginRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<LoginRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: LoginRequest): LoginRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional string username */ 1:
                    message.username = reader.string();
                    break;
                case /* optional string password */ 2:
                    message.password = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: LoginRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional string username = 1; */
        if (message.username !== undefined)
            writer.tag(1, WireType.LengthDelimited).string(message.username);
        /* optional string password = 2; */
        if (message.password !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.password);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message api_auth.LoginRequest
 */
export const LoginRequest = new LoginRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class LoginResponse$Type extends MessageType<LoginResponse> {
    constructor() {
        super("api_auth.LoginResponse", [
            { no: 1, name: "status", kind: "message", T: () => Status },
            { no: 2, name: "token", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<LoginResponse>): LoginResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<LoginResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: LoginResponse): LoginResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional common.Status status */ 1:
                    message.status = Status.internalBinaryRead(reader, reader.uint32(), options, message.status);
                    break;
                case /* optional string token */ 2:
                    message.token = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: LoginResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional common.Status status = 1; */
        if (message.status)
            Status.internalBinaryWrite(message.status, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional string token = 2; */
        if (message.token !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.token);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message api_auth.LoginResponse
 */
export const LoginResponse = new LoginResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RegisterRequest$Type extends MessageType<RegisterRequest> {
    constructor() {
        super("api_auth.RegisterRequest", [
            { no: 1, name: "username", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "password", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "email", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<RegisterRequest>): RegisterRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<RegisterRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RegisterRequest): RegisterRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional string username */ 1:
                    message.username = reader.string();
                    break;
                case /* optional string password */ 2:
                    message.password = reader.string();
                    break;
                case /* optional string email */ 3:
                    message.email = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: RegisterRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional string username = 1; */
        if (message.username !== undefined)
            writer.tag(1, WireType.LengthDelimited).string(message.username);
        /* optional string password = 2; */
        if (message.password !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.password);
        /* optional string email = 3; */
        if (message.email !== undefined)
            writer.tag(3, WireType.LengthDelimited).string(message.email);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message api_auth.RegisterRequest
 */
export const RegisterRequest = new RegisterRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RegisterResponse$Type extends MessageType<RegisterResponse> {
    constructor() {
        super("api_auth.RegisterResponse", [
            { no: 1, name: "status", kind: "message", T: () => Status }
        ]);
    }
    create(value?: PartialMessage<RegisterResponse>): RegisterResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<RegisterResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RegisterResponse): RegisterResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional common.Status status */ 1:
                    message.status = Status.internalBinaryRead(reader, reader.uint32(), options, message.status);
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: RegisterResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional common.Status status = 1; */
        if (message.status)
            Status.internalBinaryWrite(message.status, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message api_auth.RegisterResponse
 */
export const RegisterResponse = new RegisterResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class APIPermission$Type extends MessageType<APIPermission> {
    constructor() {
        super("api_auth.APIPermission", [
            { no: 1, name: "method", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "path", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<APIPermission>): APIPermission {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<APIPermission>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: APIPermission): APIPermission {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional string method */ 1:
                    message.method = reader.string();
                    break;
                case /* optional string path */ 2:
                    message.path = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: APIPermission, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional string method = 1; */
        if (message.method !== undefined)
            writer.tag(1, WireType.LengthDelimited).string(message.method);
        /* optional string path = 2; */
        if (message.path !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.path);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message api_auth.APIPermission
 */
export const APIPermission = new APIPermission$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SignTokenRequest$Type extends MessageType<SignTokenRequest> {
    constructor() {
        super("api_auth.SignTokenRequest", [
            { no: 1, name: "expires_in", kind: "scalar", opt: true, T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "permissions", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => APIPermission }
        ]);
    }
    create(value?: PartialMessage<SignTokenRequest>): SignTokenRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.permissions = [];
        if (value !== undefined)
            reflectionMergePartial<SignTokenRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SignTokenRequest): SignTokenRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional int64 expires_in */ 1:
                    message.expiresIn = reader.int64().toBigInt();
                    break;
                case /* repeated api_auth.APIPermission permissions */ 2:
                    message.permissions.push(APIPermission.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: SignTokenRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional int64 expires_in = 1; */
        if (message.expiresIn !== undefined)
            writer.tag(1, WireType.Varint).int64(message.expiresIn);
        /* repeated api_auth.APIPermission permissions = 2; */
        for (let i = 0; i < message.permissions.length; i++)
            APIPermission.internalBinaryWrite(message.permissions[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message api_auth.SignTokenRequest
 */
export const SignTokenRequest = new SignTokenRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SignTokenResponse$Type extends MessageType<SignTokenResponse> {
    constructor() {
        super("api_auth.SignTokenResponse", [
            { no: 1, name: "status", kind: "message", T: () => Status },
            { no: 2, name: "token", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<SignTokenResponse>): SignTokenResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<SignTokenResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SignTokenResponse): SignTokenResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional common.Status status */ 1:
                    message.status = Status.internalBinaryRead(reader, reader.uint32(), options, message.status);
                    break;
                case /* optional string token */ 2:
                    message.token = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: SignTokenResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional common.Status status = 1; */
        if (message.status)
            Status.internalBinaryWrite(message.status, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional string token = 2; */
        if (message.token !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.token);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message api_auth.SignTokenResponse
 */
export const SignTokenResponse = new SignTokenResponse$Type();
