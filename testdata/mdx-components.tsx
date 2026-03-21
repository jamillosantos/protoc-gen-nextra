import { useMDXComponents as getDocsMDXComponents } from 'nextra-theme-docs'
import { TypePreview } from './components/TypePreview'

export function useMDXComponents(components = {}) {
  return getDocsMDXComponents({ ...components, TypePreview })
}
