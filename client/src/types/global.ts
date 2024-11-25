export interface FileSystemResponse {
  dirs: Content[]
  files: Content[]
}

export interface Content {
  url: string
  name: string
}