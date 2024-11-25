import React from "react";
import File from "../components/icons/File.astro";

interface Resource {
  type: string
  url: string
  text: string
}

const Modal: React.FC<{ item: Resource | null; onClose: () => void }> = ({ item, onClose }) => {
  if (!item) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white p-4 rounded-lg max-w-3xl max-h-[90vh] overflow-auto" onClick={e => e.stopPropagation()}>
        {item.type === 'imagen' ? (
          <img src={item.url} alt={item.text} className="max-w-full h-auto" />
        ) : (
          <div className="flex flex-col items-center">
            <File />
            <p className="mt-4 text-lg font-semibold">{item.text}</p>
            <a href={item.url} download className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors">
              Descargar archivo
            </a>
          </div>
        )}
      </div>
    </div>
  )
}

export default Modal
