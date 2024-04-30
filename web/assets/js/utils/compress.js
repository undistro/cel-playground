/**
 * Copyright 2024 Undistro Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export function getDecompressedContent(content) {
  const decodedUint8Array = new Uint8Array(
    atob(content)
      .split("")
      .map(function (char) {
        return char.charCodeAt(0);
      })
  );

  const decompressedData = pako.ungzip(decodedUint8Array, { to: "string" });
  if (!decompressedData) {
    throw new Error("Invalid content parameter");
  }

  return JSON.parse(decompressedData);
}
