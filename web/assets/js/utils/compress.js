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
