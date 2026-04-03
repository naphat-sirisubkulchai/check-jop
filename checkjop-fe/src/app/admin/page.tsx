"use client";

import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { Filter, Search, Plus, RefreshCw, BookOpen, Calendar, CheckCircle, XCircle, Trash2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { courseApi } from "@/api/courseApi";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Curriculum } from "@/types";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";

export default function AdminHomePage() {
  const router = useRouter();

  // ---------- Curricula state ----------
  const [curricula, setCurricula] = useState<Curriculum[]>([]);
  const [curriculumQuery, setCurriculumQuery] = useState("");
  const [filterYear, setFilterYear] = useState<string>("all");
  const [curriculumLoading, setCurriculumLoading] = useState(false);
  
  // ---------- Loaders ----------
  async function loadCurricula() {
    try {
      setCurriculumLoading(true);
      const data = await courseApi.getAllCurriculum();
      setCurricula(data);
    } catch (e: unknown) {
      toast.error("Failed to load curricula", { description: e instanceof Error ? e.message : "Unknown error occurred" });
    } finally {
      setCurriculumLoading(false);
    }
  }

  useEffect(() => {
    loadCurricula();
  }, []);
  
  // ---------- Filters ----------
  const filteredCurricula = useMemo(() => {
    const q = curriculumQuery.trim().toLowerCase();
    if (!q) return curricula;
    return curricula
      .filter(
        (c) =>
          c.id.toLowerCase().includes(q) ||
          c.nameEN.toLowerCase().includes(q) ||
          c.nameTH.toLowerCase().includes(q) ||
          String(c.year).includes(q)
      )
      .filter((c) => {
        if (filterYear === "all") return true;
        return String(c.year) === filterYear;
      });
  }, [curriculumQuery, curricula, filterYear]);
  
  // Available years for filter
  const availableYears = useMemo(() => {
    return Array.from(new Set(curricula.map((c) => c.year))).sort(
      (a, b) => b - a
    );
  }, [curricula]);
  
  // ---------- Curriculum CRUD ----------

  async function deleteCurriculum(id: string, name: string) {
    toast.warning(`Delete curriculum: ${name}?`, {
      action: {
        label: "Confirm",
        onClick: async () => {
          try {
            await courseApi.deleteCurriculum(id);
            setCurricula((prev) => prev.filter((c) => c.id !== id));
            toast.success("Curriculum deleted successfully");
          } catch (e: unknown) {
            toast.error("Failed to delete curriculum", {
              description: e instanceof Error ? e.message : "Unknown error occurred"
            });
          }
        },
      },
    });
  }

  // Stats
  const stats = useMemo(() => ({
    total: curricula.length,
    active: curricula.filter(c => c.isActive).length,
    years: availableYears.length,
  }), [curricula, availableYears]);

  return (
    <div className="flex flex-col h-full overflow-hidden">
      {/* Main Content */}
      <main className="flex-1 overflow-y-auto">
        <div className="max-w-7xl mx-auto p-6 space-y-6">
          {/* Page Header */}
          <div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-chula-active to-pink-500 bg-clip-text text-transparent">
              Curriculum Management
            </h1>
            <p className="text-gray-600 mt-1">
              Manage and organize curriculum data across different years
            </p>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card className="p-5 bg-white border-l-4 border-l-chula-active shadow-md hover:shadow-lg transition-shadow">
              <div className="flex items-center gap-4">
                <div className="p-3 bg-gradient-to-br from-chula-active to-pink-500 rounded-xl">
                  <BookOpen className="h-6 w-6 text-white" />
                </div>
                <div>
                  <p className="text-sm text-gray-600 font-medium">Total Curricula</p>
                  <p className="text-3xl font-bold text-gray-900">{stats.total}</p>
                </div>
              </div>
            </Card>

            <Card className="p-5 bg-white border-l-4 border-l-green-500 shadow-md hover:shadow-lg transition-shadow">
              <div className="flex items-center gap-4">
                <div className="p-3 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl">
                  <CheckCircle className="h-6 w-6 text-white" />
                </div>
                <div>
                  <p className="text-sm text-gray-600 font-medium">Active</p>
                  <p className="text-3xl font-bold text-gray-900">{stats.active}</p>
                </div>
              </div>
            </Card>

            <Card className="p-5 bg-white border-l-4 border-l-blue-500 shadow-md hover:shadow-lg transition-shadow">
              <div className="flex items-center gap-4">
                <div className="p-3 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl">
                  <Calendar className="h-6 w-6 text-white" />
                </div>
                <div>
                  <p className="text-sm text-gray-600 font-medium">Academic Years</p>
                  <p className="text-3xl font-bold text-gray-900">{stats.years}</p>
                </div>
              </div>
            </Card>
          </div>

          {/* Search & Actions */}
          <Card className="p-2 bg-white shadow-md">
            <div className="flex flex-col md:flex-row gap-3">
              <div className="flex-1 relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search by ID, name, or year..."
                  value={curriculumQuery}
                  onChange={(e) => setCurriculumQuery(e.target.value)}
                  className="pl-9"
                />
              </div>

              <Select value={filterYear} onValueChange={setFilterYear}>
                <SelectTrigger className="w-full md:w-48">
                  <Filter className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Filter by year" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Years</SelectItem>
                  {availableYears.map((year) => (
                    <SelectItem key={year} value={year.toString()}>
                      Year {year}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>

              <Button
                onClick={loadCurricula}
                disabled={curriculumLoading}
                variant="outline"
                className="gap-2"
              >
                <RefreshCw className={`h-4 w-4 ${curriculumLoading ? "animate-spin" : ""}`} />
                Refresh
              </Button>

              <Button
                onClick={() => router.push("/admin/curriculum/create")}
                className="bg-gradient-to-r from-chula-active to-pink-500 hover:from-chula-active/90 hover:to-pink-600 gap-2"
              >
                <Plus className="h-4 w-4" />
                Add Curriculum
              </Button>
            </div>
          </Card>

          {/* Table */}
          <Card className="overflow-hidden shadow-md p-0">
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead className="bg-gradient-to-r from-gray-50 to-gray-100 border-b-2 border-gray-200">
                  <tr>
                    <th className="text-left p-4 font-semibold text-gray-700">ID</th>
                    <th className="text-left p-4 font-semibold text-gray-700">Name (TH)</th>
                    <th className="text-left p-4 font-semibold text-gray-700">Name (EN)</th>
                    <th className="text-left p-4 font-semibold text-gray-700">Year</th>
                    <th className="text-left p-4 font-semibold text-gray-700">Min Credits</th>
                    <th className="text-left p-4 font-semibold text-gray-700">Status</th>
                    <th className="text-left p-4 font-semibold text-gray-700">Actions</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-100">
                  {filteredCurricula.map((c) => (
                    <tr key={c.id} className="hover:bg-gray-50 transition-colors">
                      <td className="p-4 font-mono text-xs text-gray-600">{c.id}</td>
                      <td className="p-4 font-medium text-gray-900">{c.nameTH}</td>
                      <td className="p-4 text-gray-700">{c.nameEN}</td>
                      <td className="p-4">
                        <Badge variant="secondary" className="bg-blue-100 text-blue-700">
                          {c.year}
                        </Badge>
                      </td>
                      <td className="p-4 text-gray-700">{c.minTotalCredits}</td>
                      <td className="p-4">
                        {c.isActive ? (
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
                      </td>
                      <td className="p-4">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => deleteCurriculum(c.id, c.nameTH)}
                          className="text-red-600 hover:text-red-700 hover:bg-red-50"
                        >
                          <Trash2 className="h-4 w-4 mr-1" />
                          Delete
                        </Button>
                      </td>
                    </tr>
                  ))}

                  {!curriculumLoading && filteredCurricula.length === 0 && (
                    <tr>
                      <td colSpan={7} className="p-12 text-center">
                        <div className="flex flex-col items-center gap-3">
                          <div className="p-4 bg-gray-100 rounded-full">
                            <BookOpen className="h-8 w-8 text-gray-400" />
                          </div>
                          <div>
                            <p className="text-gray-900 font-medium">No curricula found</p>
                            <p className="text-sm text-gray-500 mt-1">
                              {curriculumQuery ? "Try adjusting your search" : "Get started by adding a curriculum"}
                            </p>
                          </div>
                        </div>
                      </td>
                    </tr>
                  )}

                  {curriculumLoading && (
                    <tr>
                      <td colSpan={7} className="p-12 text-center">
                        <div className="flex items-center justify-center gap-3">
                          <RefreshCw className="h-5 w-5 text-chula-active animate-spin" />
                          <p className="text-gray-600">Loading curricula...</p>
                        </div>
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </Card>
        </div>
      </main>
    </div>
  );
}
