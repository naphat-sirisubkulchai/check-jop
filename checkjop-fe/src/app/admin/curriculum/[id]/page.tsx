"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { courseApi } from "@/api/courseApi";
import { Curriculum, Course } from "@/types";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import {
  ArrowLeft,
  BookOpen,
  CheckCircle,
  XCircle,
  RefreshCw,
  GitBranch,
  Layers,
  GraduationCap,
} from "lucide-react";
import dynamic from "next/dynamic";
import "@xyflow/react/dist/style.css";

const CourseGraph = dynamic(() => import("@/graph/components/CourseGraph"), {
  ssr: false,
  loading: () => (
    <div className="flex items-center justify-center h-64 text-gray-500">
      <RefreshCw className="animate-spin mr-2 h-5 w-5" />
      Loading graph...
    </div>
  ),
});

export default function CurriculumDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const [curriculum, setCurriculum] = useState<Curriculum | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      setLoading(true);
      const data = await courseApi.getCurriculumById(id);
      setCurriculum(data);
      setLoading(false);
    }
    load();
  }, [id]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64 text-gray-500">
        <RefreshCw className="animate-spin mr-2 h-5 w-5" />
        Loading curriculum...
      </div>
    );
  }

  if (!curriculum) {
    return (
      <div className="max-w-7xl mx-auto p-6">
        <Button variant="ghost" onClick={() => router.back()} className="mb-4">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
        <p className="text-gray-500">Curriculum not found.</p>
      </div>
    );
  }

  // curriculum.courses already has prereqs/coreqs as strings, set by getCurriculumById
  const coursesForGraph: Course[] = curriculum.courses ?? [];

  return (
    <div className="flex flex-col h-full overflow-hidden">
      <main className="flex-1 overflow-y-auto">
        <div className="max-w-7xl mx-auto p-6 space-y-6">
          {/* Back + Header */}
          <div>
            <Button variant="ghost" onClick={() => router.back()} className="mb-2 -ml-2">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Admin
            </Button>
            <div className="flex items-start justify-between">
              <div>
                <h1 className="text-2xl font-bold text-gray-900">{curriculum.nameTH}</h1>
                <p className="text-gray-500 mt-0.5">{curriculum.nameEN}</p>
              </div>
              <div className="flex items-center gap-2">
                {curriculum.isActive ? (
                  <Badge className="bg-green-100 text-green-700 border-green-200">
                    <CheckCircle className="h-3 w-3 mr-1" />
                    Active
                  </Badge>
                ) : (
                  <Badge className="bg-gray-100 text-gray-600 border-gray-200">
                    <XCircle className="h-3 w-3 mr-1" />
                    Inactive
                  </Badge>
                )}
              </div>
            </div>
          </div>

          {/* Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <Card className="p-4 bg-white border border-gray-100 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-chula-soft rounded-lg">
                  <GraduationCap className="h-5 w-5 text-chula-active" />
                </div>
                <div>
                  <p className="text-xs text-gray-500">Academic Year</p>
                  <p className="text-xl font-bold text-gray-900">{curriculum.year}</p>
                </div>
              </div>
            </Card>
            <Card className="p-4 bg-white border border-gray-100 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-blue-50 rounded-lg">
                  <BookOpen className="h-5 w-5 text-blue-600" />
                </div>
                <div>
                  <p className="text-xs text-gray-500">Min Credits</p>
                  <p className="text-xl font-bold text-gray-900">{curriculum.minTotalCredits}</p>
                </div>
              </div>
            </Card>
            <Card className="p-4 bg-white border border-gray-100 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-purple-50 rounded-lg">
                  <Layers className="h-5 w-5 text-purple-600" />
                </div>
                <div>
                  <p className="text-xs text-gray-500">Categories</p>
                  <p className="text-xl font-bold text-gray-900">{curriculum.categories?.length ?? 0}</p>
                </div>
              </div>
            </Card>
            <Card className="p-4 bg-white border border-gray-100 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-green-50 rounded-lg">
                  <GitBranch className="h-5 w-5 text-green-600" />
                </div>
                <div>
                  <p className="text-xs text-gray-500">Courses</p>
                  <p className="text-xl font-bold text-gray-900">{coursesForGraph.length}</p>
                </div>
              </div>
            </Card>
          </div>

          {/* Tabs */}
          <Tabs defaultValue="categories">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="categories">Categories & Courses</TabsTrigger>
              <TabsTrigger value="graph">Dependency Graph</TabsTrigger>
            </TabsList>

            {/* Categories & Courses Tab */}
            <TabsContent value="categories" className="mt-4 space-y-4">
              {(!curriculum.categories || curriculum.categories.length === 0) && (
                <p className="text-sm text-gray-500">No categories found.</p>
              )}
              {[...(curriculum.categories ?? [])].sort((a: any, b: any) => b.minCredits - a.minCredits).map((cat: any) => (
                <Card key={cat.id} className="p-0 overflow-hidden border border-gray-100 shadow-sm">
                  <div className="flex items-center justify-between px-4 py-3 bg-gray-50 border-b border-gray-100">
                    <div>
                      <p className="font-semibold text-gray-900">{cat.nameTH}</p>
                      <p className="text-xs text-gray-500">{cat.nameEN}</p>
                    </div>
                    <Badge variant="secondary" className="text-xs">
                      {cat.minCredits} credits required
                    </Badge>
                  </div>
                  {cat.courses && cat.courses.length > 0 ? (
                    <div className="overflow-x-auto">
                      <table className="w-full text-sm">
                        <thead className="bg-white border-b border-gray-100">
                          <tr>
                            <th className="text-left px-4 py-2 font-medium text-gray-600 w-24">Code</th>
                            <th className="text-left px-4 py-2 font-medium text-gray-600">Name (TH)</th>
                            <th className="text-left px-4 py-2 font-medium text-gray-600">Name (EN)</th>
                            <th className="text-left px-4 py-2 font-medium text-gray-600 w-16">Credits</th>
                            <th className="text-left px-4 py-2 font-medium text-gray-600">Prerequisites</th>
                            <th className="text-left px-4 py-2 font-medium text-gray-600">Corequisites</th>
                          </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-50">
                          {cat.courses.map((course: Course) => (
                            <tr key={course.code} className="hover:bg-gray-50">
                              <td className="px-4 py-2 font-mono text-xs text-gray-700">{course.code}</td>
                              <td className="px-4 py-2 text-gray-900">{course.nameTH}</td>
                              <td className="px-4 py-2 text-gray-600">{course.nameEN}</td>
                              <td className="px-4 py-2 text-gray-700">{course.credits}</td>
                              <td className="px-4 py-2 text-xs text-gray-500 font-mono whitespace-nowrap">
                                {course.prerequisites || "—"}
                              </td>
                              <td className="px-4 py-2 text-xs text-gray-500 font-mono whitespace-nowrap">
                                {course.corequisites || "—"}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  ) : (
                    <p className="px-4 py-3 text-sm text-gray-400">No courses in this category.</p>
                  )}
                </Card>
              ))}
            </TabsContent>

            {/* Dependency Graph Tab */}
            <TabsContent value="graph" className="mt-4">
              <div className="mb-3 flex items-center gap-6 text-sm text-gray-600">
                <div className="flex items-center gap-2">
                  <div className="w-6 h-0.5 bg-red-500"></div>
                  <span>Prerequisites</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-6 h-0.5 bg-green-500 border-dashed border-2 border-green-500"></div>
                  <span>Corequisites</span>
                </div>
              </div>
              <Card className="overflow-hidden border border-gray-100 shadow-sm" style={{ height: "65vh" }}>
                {coursesForGraph.length > 0 ? (
                  <CourseGraph courses={coursesForGraph} />
                ) : (
                  <div className="flex items-center justify-center h-full text-gray-400">
                    No courses with dependency data.
                  </div>
                )}
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </main>
    </div>
  );
}
