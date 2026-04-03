'use client';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Home, Settings, Beaker, BookOpen } from "lucide-react"

export default function Portal() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100 p-4">
      <Card className="w-full max-w-5xl shadow-lg rounded-2xl border-0 bg-white">
        <CardHeader className="pb-6 border-b">
          <CardTitle className="text-3xl font-bold text-gray-900 text-center">
            CheckJop Portal
          </CardTitle>
          <p className="text-center text-gray-600 mt-2">
            Select your destination
          </p>
        </CardHeader>
        <CardContent className="pt-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Setup Page Link */}
          <Link href="/setup" className="group">
            <div className="h-full p-6 border-2 border-gray-200 rounded-xl hover:border-chula-active hover:bg-chula-soft/30 transition-all duration-200 cursor-pointer">
              <div className="flex flex-col items-center justify-center space-y-4">
                <div className="w-16 h-16 rounded-full bg-gradient-to-br from-chula-soft to-pink-100 flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                  <BookOpen className="w-8 h-8 text-chula-active" />
                </div>
                <div className="text-center">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">
                    Setup
                  </h3>
                  <p className="text-sm text-gray-600">
                    Configure curriculum & start planning
                  </p>
                </div>
              </div>
            </div>
          </Link>

          {/* Home Page Link */}
          <Link href="/home" className="group">
            <div className="h-full p-6 border-2 border-gray-200 rounded-xl hover:border-blue-400 hover:bg-blue-50/30 transition-all duration-200 cursor-pointer">
              <div className="flex flex-col items-center justify-center space-y-4">
                <div className="w-16 h-16 rounded-full bg-gradient-to-br from-blue-100 to-indigo-100 flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                  <Home className="w-8 h-8 text-blue-600" />
                </div>
                <div className="text-center">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">
                    Study Plan
                  </h3>
                  <p className="text-sm text-gray-600">
                    View and manage your study plan
                  </p>
                </div>
              </div>
            </div>
          </Link>

          {/* Admin Page Link */}
          <Link href="/admin" className="group">
            <div className="h-full p-6 border-2 border-gray-200 rounded-xl hover:border-purple-400 hover:bg-purple-50/30 transition-all duration-200 cursor-pointer">
              <div className="flex flex-col items-center justify-center space-y-4">
                <div className="w-16 h-16 rounded-full bg-gradient-to-br from-purple-100 to-pink-100 flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                  <Settings className="w-8 h-8 text-purple-600" />
                </div>
                <div className="text-center">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">
                    Admin
                  </h3>
                  <p className="text-sm text-gray-600">
                    Manage curriculum and courses
                  </p>
                </div>
              </div>
            </div>
          </Link>

          {/* Testing Page Link */}
          <Link href="/testing" className="group">
            <div className="h-full p-6 border-2 border-gray-200 rounded-xl hover:border-green-400 hover:bg-green-50/30 transition-all duration-200 cursor-pointer">
              <div className="flex flex-col items-center justify-center space-y-4">
                <div className="w-16 h-16 rounded-full bg-gradient-to-br from-green-100 to-emerald-100 flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                  <Beaker className="w-8 h-8 text-green-600" />
                </div>
                <div className="text-center">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">
                    Testing
                  </h3>
                  <p className="text-sm text-gray-600">
                    Load test cases and validate
                  </p>
                </div>
              </div>
            </div>
          </Link>
        </CardContent>
      </Card>
    </div>
  )
}
