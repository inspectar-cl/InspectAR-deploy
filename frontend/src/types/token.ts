export interface DecodedJwt {
    username: string;
    device: string;
    exp: number;
    scope: string;
}