'use client';

import { useState } from 'react';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { GraduationCap, CheckCircle, Calendar, BookOpen, ArrowRight, Loader2 } from 'lucide-react';

export default function LandingPage() {
  const router = useRouter();
  const [isLoggingIn, setIsLoggingIn] = useState(false);

  const handleSSOLogin = async () => {
    setIsLoggingIn(true);

    try {
      router.push('/setup');
    } catch (error) {
      console.error('Login failed:', error);
    } finally {
      setIsLoggingIn(false);
    }
  };

  const features = [
    {
      icon: <CheckCircle className="h-6 w-6 text-chula-active" />,
      title: 'Check Graduation',
      description: 'Verify your graduation eligibility instantly',
    },
    {
      icon: <Calendar className="h-6 w-6 text-chula-active" />,
      title: 'Plan Your Journey',
      description: 'Organize courses by semester and year',
    },
    {
      icon: <BookOpen className="h-6 w-6 text-chula-active" />,
      title: 'Track Progress',
      description: 'Monitor credits and requirements in real-time',
    },
  ];

  return (
    <div>
      {/* Hero Section */}
      <main className="container mx-auto px-6 py-12 md:py-20">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left Column - Hero Content */}
          <div className="space-y-8">
            <div className="space-y-4">
              <div className="inline-flex items-center gap-2 px-4 py-2 bg-white rounded-full shadow-sm border border-gray-200">
                <GraduationCap className="h-5 w-5 text-chula-active" />
                <span className="text-sm font-medium text-gray-700">
                  Chulalongkorn University
                </span>
              </div>

              <h2 className="text-5xl md:text-6xl font-bold text-gray-900 leading-tight">
                Plan Your Path to{' '}
                <span className="bg-gradient-to-r from-chula-active to-pink-500 bg-clip-text text-transparent">
                  Graduation
                </span>
              </h2>

              <p className="text-xl text-gray-600 leading-relaxed">
                The smart way to organize your study plan, track your progress, and ensure you meet all graduation requirements.
              </p>
            </div>

            {/* CTA Button */}
            <div className="space-y-4">
              <Button
                onClick={handleSSOLogin}
                disabled={isLoggingIn}
                className="h-14 px-8 bg-gradient-to-r from-chula-active to-pink-500 hover:from-chula-active/90 hover:to-pink-600 text-white text-lg font-semibold shadow-lg hover:shadow-xl transition-all duration-200 hover:scale-[1.02] active:scale-[0.98]"
                size="lg"
              >
                {isLoggingIn ? (
                  <>
                    <Loader2 className="h-5 w-5 mr-2 animate-spin" />
                    Signing In...
                  </>
                ) : (
                  <>
                    Let's go
                    <ArrowRight className="h-5 w-5 ml-2" />
                  </>
                )}
              </Button>

              <p className="text-sm text-gray-500">
                Use your Chula account to get started
              </p>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-3 gap-6 pt-8 border-t border-gray-200">
              <div>
                <div className="text-3xl font-bold text-gray-900">10K+</div>
                <div className="text-sm text-gray-600">Students</div>
              </div>
              <div>
                <div className="text-3xl font-bold text-gray-900">50+</div>
                <div className="text-sm text-gray-600">Curriculums</div>
              </div>
              <div>
                <div className="text-3xl font-bold text-gray-900">99%</div>
                <div className="text-sm text-gray-600">Accuracy</div>
              </div>
            </div>
          </div>

          {/* Right Column - Features Cards */}
          <div className="space-y-6">
            <Card className="p-8 bg-white/80 backdrop-blur-sm border-2 border-gray-100 shadow-xl">
              <div className="space-y-6">
                <div className="flex items-center gap-3 mb-6">
                  <div className="p-3 bg-gradient-to-br from-chula-active to-pink-500 rounded-xl shadow-lg">
                    <GraduationCap className="h-6 w-6 text-white" />
                  </div>
                  <h3 className="text-2xl font-bold text-gray-900">
                    Key Features
                  </h3>
                </div>

                {features.map((feature, index) => (
                  <div
                    key={index}
                    className="flex gap-4 p-4 rounded-xl hover:bg-gray-50 transition-colors"
                  >
                    <div className="flex size-10 items-center justify-center bg-chula-soft rounded-lg">
                      {feature.icon}
                    </div>
                    <div>
                      <h4 className="font-semibold text-gray-900 mb-1">
                        {feature.title}
                      </h4>
                      <p className="text-sm text-gray-600">
                        {feature.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </Card>

            {/* Quick Access Card */}
            <Card className="p-6 bg-gradient-to-r from-chula-active to-pink-500 text-white shadow-xl">
              <div className="space-y-3">
                <h4 className="font-semibold text-lg">Already have an account?</h4>
                <p className="text-sm text-white/90">
                  Sign in to access your saved study plans and continue planning your academic journey.
                </p>
              </div>
            </Card>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="container mx-auto px-6 py-8 mt-20 border-t border-gray-200">
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-sm text-gray-600">
            © 2024 CheckJop. All rights reserved.
          </p>
          <div className="flex gap-6 text-sm text-gray-600">
            <button className="hover:text-chula-active transition-colors">
              Privacy Policy
            </button>
            <button className="hover:text-chula-active transition-colors">
              Terms of Service
            </button>
            <button className="hover:text-chula-active transition-colors">
              Contact Us
            </button>
          </div>
        </div>
      </footer>
    </div>
  );
}
