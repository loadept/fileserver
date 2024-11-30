import type { APIRoute } from "astro"
import apiClient from "../../../config/apiClient"


export const POST: APIRoute = async ({ params, request }) => {
  try {
    const directory = params.directory || ''

    const data = await request.formData()
    
    const cookies = request.headers.get('cookie') || '';
    const token = cookies.split('=')[1]

    const { data: response } = await apiClient.put(`fs/upload?directory=${directory}`, data, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${token}`
      }
    })

    console.log(response)

    return new Response(
      JSON.stringify({
        message: 'Success'
      }),
      { status: 200 }
    )
  } catch (err) {
    console.error((err as Error).message)
    return new Response(
      JSON.stringify({
        error: 'Internal server error'
      }),
      { status: 500 }
    )
  }
}
