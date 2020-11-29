package binarypack

// func Pack(data []byte) []byte {
//     packer := newPacker()
//     packer.pack(data)
//     buffer := packer.getBuffer()
//     return buffer;
// }

// function Packer () {
//   this.bufferBuilder = new BufferBuilder();
// }

// Packer.prototype.getBuffer = function () {
//   return this.bufferBuilder.getBuffer();
// };

// Packer.prototype.pack = function (value) {
//   var type = typeof (value);
//   if (type === 'string') {
//     this.pack_string(value);
//   } else if (type === 'number') {
//     if (Math.floor(value) === value) {
//       this.pack_integer(value);
//     } else {
//       this.pack_double(value);
//     }
//   } else if (type === 'boolean') {
//     if (value === true) {
//       this.bufferBuilder.append(0xc3);
//     } else if (value === false) {
//       this.bufferBuilder.append(0xc2);
//     }
//   } else if (type === 'undefined') {
//     this.bufferBuilder.append(0xc0);
//   } else if (type === 'object') {
//     if (value === null) {
//       this.bufferBuilder.append(0xc0);
//     } else {
//       var constructor = value.constructor;
//       if (constructor == Array) {
//         this.pack_array(value);
//       } else if (constructor == Blob || constructor == File || value instanceof Blob || value instanceof File) {
//         this.pack_bin(value);
//       } else if (constructor == ArrayBuffer) {
//         if (binaryFeatures.useArrayBufferView) {
//           this.pack_bin(new Uint8Array(value));
//         } else {
//           this.pack_bin(value);
//         }
//       } else if ('BYTES_PER_ELEMENT' in value) {
//         if (binaryFeatures.useArrayBufferView) {
//           this.pack_bin(new Uint8Array(value.buffer));
//         } else {
//           this.pack_bin(value.buffer);
//         }
//       } else if ((constructor == Object) || (constructor.toString().startsWith('class'))) {
//         this.pack_object(value);
//       } else if (constructor == Date) {
//         this.pack_string(value.toString());
//       } else if (typeof value.toBinaryPack === 'function') {
//         this.bufferBuilder.append(value.toBinaryPack());
//       } else {
//         throw new Error('Type "' + constructor.toString() + '" not yet supported');
//       }
//     }
//   } else {
//     throw new Error('Type "' + type + '" not yet supported');
//   }
//   this.bufferBuilder.flush();
// };

// Packer.prototype.pack_bin = function (blob) {
//   var length = blob.length || blob.byteLength || blob.size;
//   if (length <= 0x0f) {
//     this.pack_uint8(0xa0 + length);
//   } else if (length <= 0xffff) {
//     this.bufferBuilder.append(0xda);
//     this.pack_uint16(length);
//   } else if (length <= 0xffffffff) {
//     this.bufferBuilder.append(0xdb);
//     this.pack_uint32(length);
//   } else {
//     throw new Error('Invalid length');
//   }
//   this.bufferBuilder.append(blob);
// };

// Packer.prototype.pack_string = function (str) {
//   var length = utf8Length(str);

//   if (length <= 0x0f) {
//     this.pack_uint8(0xb0 + length);
//   } else if (length <= 0xffff) {
//     this.bufferBuilder.append(0xd8);
//     this.pack_uint16(length);
//   } else if (length <= 0xffffffff) {
//     this.bufferBuilder.append(0xd9);
//     this.pack_uint32(length);
//   } else {
//     throw new Error('Invalid length');
//   }
//   this.bufferBuilder.append(str);
// };

// Packer.prototype.pack_array = function (ary) {
//   var length = ary.length;
//   if (length <= 0x0f) {
//     this.pack_uint8(0x90 + length);
//   } else if (length <= 0xffff) {
//     this.bufferBuilder.append(0xdc);
//     this.pack_uint16(length);
//   } else if (length <= 0xffffffff) {
//     this.bufferBuilder.append(0xdd);
//     this.pack_uint32(length);
//   } else {
//     throw new Error('Invalid length');
//   }
//   for (var i = 0; i < length; i++) {
//     this.pack(ary[i]);
//   }
// };

// Packer.prototype.pack_integer = function (num) {
//   if (num >= -0x20 && num <= 0x7f) {
//     this.bufferBuilder.append(num & 0xff);
//   } else if (num >= 0x00 && num <= 0xff) {
//     this.bufferBuilder.append(0xcc);
//     this.pack_uint8(num);
//   } else if (num >= -0x80 && num <= 0x7f) {
//     this.bufferBuilder.append(0xd0);
//     this.pack_int8(num);
//   } else if (num >= 0x0000 && num <= 0xffff) {
//     this.bufferBuilder.append(0xcd);
//     this.pack_uint16(num);
//   } else if (num >= -0x8000 && num <= 0x7fff) {
//     this.bufferBuilder.append(0xd1);
//     this.pack_int16(num);
//   } else if (num >= 0x00000000 && num <= 0xffffffff) {
//     this.bufferBuilder.append(0xce);
//     this.pack_uint32(num);
//   } else if (num >= -0x80000000 && num <= 0x7fffffff) {
//     this.bufferBuilder.append(0xd2);
//     this.pack_int32(num);
//   } else if (num >= -0x8000000000000000 && num <= 0x7FFFFFFFFFFFFFFF) {
//     this.bufferBuilder.append(0xd3);
//     this.pack_int64(num);
//   } else if (num >= 0x0000000000000000 && num <= 0xFFFFFFFFFFFFFFFF) {
//     this.bufferBuilder.append(0xcf);
//     this.pack_uint64(num);
//   } else {
//     throw new Error('Invalid integer');
//   }
// };

// Packer.prototype.pack_double = function (num) {
//   var sign = 0;
//   if (num < 0) {
//     sign = 1;
//     num = -num;
//   }
//   var exp = Math.floor(Math.log(num) / Math.LN2);
//   var frac0 = num / Math.pow(2, exp) - 1;
//   var frac1 = Math.floor(frac0 * Math.pow(2, 52));
//   var b32 = Math.pow(2, 32);
//   var h32 = (sign << 31) | ((exp + 1023) << 20) |
//     (frac1 / b32) & 0x0fffff;
//   var l32 = frac1 % b32;
//   this.bufferBuilder.append(0xcb);
//   this.pack_int32(h32);
//   this.pack_int32(l32);
// };

// Packer.prototype.pack_object = function (obj) {
//   var keys = Object.keys(obj);
//   var length = keys.length;
//   if (length <= 0x0f) {
//     this.pack_uint8(0x80 + length);
//   } else if (length <= 0xffff) {
//     this.bufferBuilder.append(0xde);
//     this.pack_uint16(length);
//   } else if (length <= 0xffffffff) {
//     this.bufferBuilder.append(0xdf);
//     this.pack_uint32(length);
//   } else {
//     throw new Error('Invalid length');
//   }
//   for (var prop in obj) {
//     if (obj.hasOwnProperty(prop)) {
//       this.pack(prop);
//       this.pack(obj[prop]);
//     }
//   }
// };

// Packer.prototype.pack_uint8 = function (num) {
//   this.bufferBuilder.append(num);
// };

// Packer.prototype.pack_uint16 = function (num) {
//   this.bufferBuilder.append(num >> 8);
//   this.bufferBuilder.append(num & 0xff);
// };

// Packer.prototype.pack_uint32 = function (num) {
//   var n = num & 0xffffffff;
//   this.bufferBuilder.append((n & 0xff000000) >>> 24);
//   this.bufferBuilder.append((n & 0x00ff0000) >>> 16);
//   this.bufferBuilder.append((n & 0x0000ff00) >>> 8);
//   this.bufferBuilder.append((n & 0x000000ff));
// };

// Packer.prototype.pack_uint64 = function (num) {
//   var high = num / Math.pow(2, 32);
//   var low = num % Math.pow(2, 32);
//   this.bufferBuilder.append((high & 0xff000000) >>> 24);
//   this.bufferBuilder.append((high & 0x00ff0000) >>> 16);
//   this.bufferBuilder.append((high & 0x0000ff00) >>> 8);
//   this.bufferBuilder.append((high & 0x000000ff));
//   this.bufferBuilder.append((low & 0xff000000) >>> 24);
//   this.bufferBuilder.append((low & 0x00ff0000) >>> 16);
//   this.bufferBuilder.append((low & 0x0000ff00) >>> 8);
//   this.bufferBuilder.append((low & 0x000000ff));
// };

// Packer.prototype.pack_int8 = function (num) {
//   this.bufferBuilder.append(num & 0xff);
// };

// Packer.prototype.pack_int16 = function (num) {
//   this.bufferBuilder.append((num & 0xff00) >> 8);
//   this.bufferBuilder.append(num & 0xff);
// };

// Packer.prototype.pack_int32 = function (num) {
//   this.bufferBuilder.append((num >>> 24) & 0xff);
//   this.bufferBuilder.append((num & 0x00ff0000) >>> 16);
//   this.bufferBuilder.append((num & 0x0000ff00) >>> 8);
//   this.bufferBuilder.append((num & 0x000000ff));
// };

// Packer.prototype.pack_int64 = function (num) {
//   var high = Math.floor(num / Math.pow(2, 32));
//   var low = num % Math.pow(2, 32);
//   this.bufferBuilder.append((high & 0xff000000) >>> 24);
//   this.bufferBuilder.append((high & 0x00ff0000) >>> 16);
//   this.bufferBuilder.append((high & 0x0000ff00) >>> 8);
//   this.bufferBuilder.append((high & 0x000000ff));
//   this.bufferBuilder.append((low & 0xff000000) >>> 24);
//   this.bufferBuilder.append((low & 0x00ff0000) >>> 16);
//   this.bufferBuilder.append((low & 0x0000ff00) >>> 8);
//   this.bufferBuilder.append((low & 0x000000ff));
// };

// function _utf8Replace (m) {
//   var code = m.charCodeAt(0);

//   if (code <= 0x7ff) return '00';
//   if (code <= 0xffff) return '000';
//   if (code <= 0x1fffff) return '0000';
//   if (code <= 0x3ffffff) return '00000';
//   return '000000';
// }

// function utf8Length (str) {
//   if (str.length > 600) {
//     // Blob method faster for large strings
//     return (new Blob([str])).size;
//   } else {
//     return str.replace(/[^\u0000-\u007F]/g, _utf8Replace).length;
//   }
// }
