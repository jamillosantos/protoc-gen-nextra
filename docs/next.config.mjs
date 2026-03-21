import nextra from 'nextra'

const isExport = process.env.NEXT_OUTPUT === 'export'
const basePath = process.env.BASE_PATH ?? ''

const withNextra = nextra({
  contentDirBasePath: '/',
})

export default withNextra({
  output: isExport ? 'export' : undefined,
  basePath,
  assetPrefix: basePath,
  images: isExport ? { unoptimized: true } : undefined,
})
