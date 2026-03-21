interface PreviewItem {
  name: string
  type?: string
  optional?: boolean
  repeated?: boolean
}

interface TypePreviewProps {
  name: string
  href: string
  repeated?: boolean
  fields: PreviewItem[]
}

const MAX_PREVIEW_ITEMS = 5

export function TypePreview({ name, href, repeated, fields }: TypePreviewProps) {
  const visible = fields.slice(0, MAX_PREVIEW_ITEMS)
  const overflow = fields.length - visible.length

  return (
    <span className="type-preview">
      <a href={href} className="x-type-link">
        <code>{name}{repeated ? '[]' : ''}</code>
      </a>

      <div className="type-preview-card">
        <div className="type-preview-header">
          <span className="type-preview-title">{name}</span>
        </div>

        <table className="type-preview-table">
          <tbody>
            {visible.map((item) => (
              <tr key={item.name} className="type-preview-row">
                <td className="type-preview-field-name">
                  {item.name}
                  {item.optional && <span className="type-preview-optional-badge">optional</span>}
                </td>
                {item.type && (
                  <td className="type-preview-field-type">
                    {item.type}{item.repeated ? '[]' : ''}
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>

        {overflow > 0 && (
          <div className="type-preview-overflow">+{overflow} more</div>
        )}
      </div>
    </span>
  )
}
