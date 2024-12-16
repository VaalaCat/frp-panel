import { ClientVersion } from "@/lib/pb/api_master"

export function NeedUpgrade(version: ClientVersion | undefined) {
    if (!(version)) return false
    if (!version.gitVersion) return false
    const versionString = version?.gitVersion
    const [a, b, c] = versionString.split('.')
    return Number(b) < 1
}