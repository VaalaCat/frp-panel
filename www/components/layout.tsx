import { Toaster } from "./ui/toaster"

export const RootLayout = ({ children }: { children: React.ReactNode }) => {
    return (
        <>
            {children}
            <Toaster />
        </>
    )
}