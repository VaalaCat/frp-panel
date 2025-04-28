import { ClientVersion } from "@/lib/pb/api_master"

export function NeedUpgrade(version: ClientVersion | undefined) {
    if (!(version)) return false
    if (!version.gitVersion) return false
    const versionString = version?.gitVersion
    const [a, b, c] = versionString.split('.')
    if (Number(b) < 1) {
        return true
    }

    console.log(Number(a), Number(b), Number(c))

    if (a=='v0' && Number(b)<=1 && Number(c) <= 10) {
        return true
    }

    return false
}