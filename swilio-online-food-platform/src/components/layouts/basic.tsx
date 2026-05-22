import React from "react";

export function BasicLayout({ children }: React.PropsWithChildren<{}>) {
    return (
        <div className="bg-background relative min-h-screen flex flex-col font-sans antialiased">
            <main className="flex-1 w-full animate-in fade-in duration-500 ease-in-out">
                {children}
            </main>
        </div>
    );
}