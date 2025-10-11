/**
 * @param token
 * @returns
 */

export function decodeJwtToken(token: string | undefined) {
  if (!token || typeof token === 'undefined') {
    return null;
  }

  try {
    //Divide el token en Header, Payload y Signature
    const parts = token.split('.');
    if (parts.length !== 3) {
      return null;
    }

    const payloadBase64 = parts[1];

    const base64 = payloadBase64.replace(/-/g, '+').replace(/_/g, '/');

    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(function(c) {
          return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        })
        .join('')
    );

    return JSON.parse(jsonPayload);
  } catch (error) {
    return null;
  }
}