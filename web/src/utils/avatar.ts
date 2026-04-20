const ALLOWED_MIME = new Set(["image/png", "image/jpeg", "image/webp"]);
export const MAX_AVATAR_BYTES = 256 * 1024;

export function isTrustedAvatarDataUrl(value?: string | null): boolean {
  if (!value) return false;
  return /^data:image\/(png|jpeg|jpg|webp);base64,[a-zA-Z0-9+/=]+$/.test(value);
}

export function validateAvatarFile(file: File): string | null {
  if (!ALLOWED_MIME.has(file.type)) {
    return "Formato inválido. Use PNG, JPEG ou WEBP.";
  }
  if (file.size <= 0 || file.size > MAX_AVATAR_BYTES) {
    return "Avatar muito grande. Limite de 256KB.";
  }
  return null;
}

export function hasAllowedMagicBytes(buffer: ArrayBuffer): boolean {
  const b = new Uint8Array(buffer);
  const isPNG =
    b.length >= 8 &&
    b[0] === 0x89 &&
    b[1] === 0x50 &&
    b[2] === 0x4e &&
    b[3] === 0x47 &&
    b[4] === 0x0d &&
    b[5] === 0x0a &&
    b[6] === 0x1a &&
    b[7] === 0x0a;
  const isJPEG = b.length >= 3 && b[0] === 0xff && b[1] === 0xd8 && b[2] === 0xff;
  const isWEBP =
    b.length >= 12 &&
    b[0] === 0x52 &&
    b[1] === 0x49 &&
    b[2] === 0x46 &&
    b[3] === 0x46 &&
    b[8] === 0x57 &&
    b[9] === 0x45 &&
    b[10] === 0x42 &&
    b[11] === 0x50;
  return isPNG || isJPEG || isWEBP;
}

