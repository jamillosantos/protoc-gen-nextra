import { Footer, Layout, Navbar } from 'nextra-theme-docs'
import { Head } from 'nextra/components'
import { getPageMap } from 'nextra/page-map'
import 'nextra-theme-docs/style.css'
import './globals.css'
import type { Metadata, Viewport } from 'next'
import type { ReactNode } from 'react'

export const metadata: Metadata = {
  description: 'gRPC service documentation',
}

export const viewport: Viewport = {
  themeColor: '#ffffff',
}

export default async function RootLayout({ children }: { children: ReactNode }) {
  const pageMap = await getPageMap()
  const navbar = (
    <Navbar
      logo={<span style={{ fontWeight: 700 }}>gRPC Docs</span>}
    />
  )
  return (
    <html lang="en" dir="ltr" suppressHydrationWarning>
      <Head />
      <body>
        <Layout
          navbar={navbar}
          footer={<Footer />}
          pageMap={pageMap}
          copyPageButton={false}
          editLink={null}
          feedback={{ content: null }}
          docsRepositoryBase="https://github.com/jamillosantos/protoc-gen-nextra"
        >
          {children}
        </Layout>
      </body>
    </html>
  )
}
