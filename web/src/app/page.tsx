import { ScrapeForm } from '@/components/ScrapeForm'

export default function Home() {
  return (
    <main className="container mx-auto px-4 py-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-2">PageMail</h1>
        <p className="text-gray-600 text-center mb-8">
          抓取网页内容并发送到您的邮箱
        </p>
        
        <ScrapeForm />
        
        <div className="mt-12 text-center text-sm text-gray-500">
          <p>未登录用户每日限制 1 次请求</p>
          <p className="mt-2">
            <a href="/auth/login" className="text-blue-500 hover:underline">登录</a>
            {' | '}
            <a href="/auth/register" className="text-blue-500 hover:underline">注册</a>
          </p>
        </div>
      </div>
    </main>
  )
}