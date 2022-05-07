import {fileTypeFromBuffer} from 'file-type'

export const handler = async (event) => {
  const request = event.Records[0].cf.request
  if (!request || !request.body || !request.body.data) {
    return {status: 400, body: 'Invalid request'}
  }

  const extension = getExtension(request)
  if (!extension) {
    return {status: 400, body: 'Invalid headers'}
  }

  const data = request.body.data.substring(0, 100)
  const buff = Buffer.from(data, 'base64')
  const type = await fileTypeFromBuffer(buff)

  if (!type || !type.ext) {
    return {status: 400, body: 'Invalid file'}
  }
  if (extension !== type.ext) {
    return {status: 400, body: 'Invalid provided file type'}
  }

  switch (type.ext) {
    case "gif":
    case "heic":
    case "jpg":
    case "png":
    case "webp":
      break
    default:
      return {status: 400, body: 'Invalid file type'}
  }
  return request
}

// extracts file extension based on Content-Type header
function getExtension(request) {
  const header = request.headers["content-type"]
  if (!header || header.length !== 1)
    return null

  const contentTypeParts = header[0].value.split("/")
  if (contentTypeParts.length !== 2)
    return null
  return contentTypeParts[1] === "jpeg" ? "jpg" : contentTypeParts[1]
}
