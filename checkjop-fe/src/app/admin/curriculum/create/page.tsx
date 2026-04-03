"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { ArrowLeft, Upload, Download, FileText, Plus, X } from "lucide-react";
import Papa from "papaparse";
import { PreviewCard } from "./components/previewCard";
import { courseApi } from "@/api/courseApi";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { CreateCurriculumForm } from "@/types";

/* eslint-disable @typescript-eslint/no-explicit-any */

// CSV transform helpers (inline)
function adaptCurriculumData(rawData: any[]) {
  return rawData.map((item) => ({
    nameTH: item.curriculumNameTH || "",
    nameEN: item.curriculumNameEN || "",
    year: parseInt(item.year || item.Year || "0"),
    minTotalCredits: parseInt(item.minTotalCredit || "0"),
    isActive: item.isActive === "true" || item.isActive === "TRUE" || true,
  }));
}

function adaptCategoryData(rawData: any[]) {
  return rawData.map((item) => ({
    nameTH: item.categoryNameTH || "",
    nameEN: item.categoryNameEn || item.categoryNameEN || "",
    curriculumName: item.curriculumName || item.cirrculumName || "",
    minCredits: parseInt(item.minCredit || "0"),
    year: parseInt(item.year || item.Year || "0"),
  }));
}

function adaptCourseData(rawData: any[]) {
  return rawData.map((item) => ({
    code: item.code || "",
    nameTH: item.courseNameTH || "",
    nameEN: item.courseNameEN || "",
    credits: parseInt(item.credit || "0"),
    prerequisites: item.prerequisites || "",
    corequisites: item.corequisites || "",
    categories: item.category || "",
    curriculum: item.curriculum || "",
    year: parseInt(item.year || item.Year || "0"),
  }));
}

const templates = [
  {
    name: "Curriculum Template",
    file: "/templates/template_checkjop_curriculum.csv",
  },
  {
    name: "Category Template",
    file: "/templates/template_checkjop_catagory.csv",
  },
  { name: "Course Template", file: "/templates/template_checkjop_course.csv" },
];

interface ParsedData {
  curriculums: any[];
  categories: any[];
  courses: any[];
}

export default function CreateCurriculumPage() {
  const router = useRouter();
  const [form, setForm] = useState<CreateCurriculumForm>({
    curriculumFile: null,
    categoryFile: null,
    courseFiles: [],
    previewData: null,
  });
  const [loading, setLoading] = useState(false);
  const [isDragOver, setIsDragOver] = useState<
    "curriculum" | "category" | undefined
  >();
  const [newCourseYear, setNewCourseYear] = useState<string>("");

  async function parseCSV<T>(file: File): Promise<T[]> {
    return new Promise((resolve, reject) => {
      Papa.parse<T>(file, {
        header: true,
        skipEmptyLines: true,
        complete: (result) => resolve(result.data),
        error: (error) => reject(error),
      });
    });
  }

  const downloadTemplate = (templateFile: string, templateName: string) => {
    const link = document.createElement("a");
    link.href = templateFile;
    link.download = templateFile.split("/").pop() || "template.csv";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    toast.success(`${templateName} downloaded successfully`);
  };

  async function readCurriculumFiles(files: {
    curriculum: File;
    category: File;
    courseFiles: Array<{ file: File; year: number }>;
  }): Promise<ParsedData> {
    // Parse raw CSV data first
    const [rawCurriculums, rawCategories] = await Promise.all([
      parseCSV<any>(files.curriculum),
      parseCSV<any>(files.category),
    ]);

    // Parse all course files and merge them
    const allRawCourses: any[] = [];
    for (const courseFile of files.courseFiles) {
      const rawCourses = await parseCSV<any>(courseFile.file);
      // Add year to each course
      const coursesWithYear = rawCourses.map((course) => ({
        ...course,
        year: courseFile.year,
        Year: courseFile.year,
      }));
      allRawCourses.push(...coursesWithYear);
    }

    // Transform data using adapters to match our types
    const curriculums = adaptCurriculumData(rawCurriculums);
    const categories = adaptCategoryData(rawCategories);
    const courses = adaptCourseData(allRawCourses);

    return { curriculums, categories, courses };
  }

  const handleFileUpload = async (
    type: "curriculum" | "category",
    file: File | null,
  ) => {
    // อัปเดต state แล้วหา currentFiles หลังจาก update
    setForm((prev) => {
      const newForm = {
        ...prev,
        [`${type}File`]: file,
      };

      setTimeout(async () => {
        const curriculumFile =
          type === "curriculum" ? file : newForm.curriculumFile;
        const categoryFile = type === "category" ? file : newForm.categoryFile;

        if (curriculumFile && categoryFile && newForm.courseFiles.length > 0) {
          try {
            const parsed = await readCurriculumFiles({
              curriculum: curriculumFile,
              category: categoryFile,
              courseFiles: newForm.courseFiles,
            });

            setForm((prevForm) => ({
              ...prevForm,
              previewData: {
                curriculums: parsed.curriculums,
                categories: parsed.categories,
                courses: parsed.courses,
                totalCurriculums: parsed.curriculums.length,
                totalCategories: parsed.categories.length,
                totalCourses: parsed.courses.length,
              },
            }));
            console.log("Parsed Data:", parsed);
          } catch (error) {
            console.error(`Error parsing CSV files:`, error);
            toast.error(`Failed to parse CSV files`);
          }
        }
      }, 0);

      return newForm;
    });
  };

  const handleAddCourseFile = (file: File, year: number) => {
    setForm((prev) => {
      const newForm = {
        ...prev,
        courseFiles: [...prev.courseFiles, { file, year }],
      };

      setTimeout(async () => {
        if (newForm.curriculumFile && newForm.categoryFile && newForm.courseFiles.length > 0) {
          try {
            const parsed = await readCurriculumFiles({
              curriculum: newForm.curriculumFile,
              category: newForm.categoryFile,
              courseFiles: newForm.courseFiles,
            });

            setForm((prevForm) => ({
              ...prevForm,
              previewData: {
                curriculums: parsed.curriculums,
                categories: parsed.categories,
                courses: parsed.courses,
                totalCurriculums: parsed.curriculums.length,
                totalCategories: parsed.categories.length,
                totalCourses: parsed.courses.length,
              },
            }));
          } catch (error) {
            console.error(`Error parsing CSV files:`, error);
            toast.error(`Failed to parse CSV files`);
          }
        }
      }, 0);

      return newForm;
    });
    setNewCourseYear("");
  };

  const handleRemoveCourseFile = (index: number) => {
    setForm((prev) => {
      const newForm = {
        ...prev,
        courseFiles: prev.courseFiles.filter((_, i) => i !== index),
      };

      setTimeout(async () => {
        if (newForm.curriculumFile && newForm.categoryFile && newForm.courseFiles.length > 0) {
          try {
            const parsed = await readCurriculumFiles({
              curriculum: newForm.curriculumFile,
              category: newForm.categoryFile,
              courseFiles: newForm.courseFiles,
            });

            setForm((prevForm) => ({
              ...prevForm,
              previewData: {
                curriculums: parsed.curriculums,
                categories: parsed.categories,
                courses: parsed.courses,
                totalCurriculums: parsed.curriculums.length,
                totalCategories: parsed.categories.length,
                totalCourses: parsed.courses.length,
              },
            }));
          } catch (error) {
            console.error(`Error parsing CSV files:`, error);
            toast.error(`Failed to parse CSV files`);
          }
        } else {
          setForm((prevForm) => ({
            ...prevForm,
            previewData: null,
          }));
        }
      }, 0);

      return newForm;
    });
  };

  const handleSubmitCurriculums = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!form.curriculumFile || !form.categoryFile || form.courseFiles.length === 0) {
      toast.warning("Please upload all required files");
      return;
    }

    try {
      setLoading(true);
      await courseApi.createCurriculum({
        fileCurriculum: form.curriculumFile,
        fileCategory: form.categoryFile,
        courseFiles: form.courseFiles,
      });
      toast.success("Curriculum created successfully");
      router.push("/admin");
    } catch (e: unknown) {
      toast.error("Failed to create curriculum", {
        description: e instanceof Error ? e.message : "Unknown error occurred",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleDragOver = (
    e: React.DragEvent,
    fileType: "curriculum" | "category",
  ) => {
    e.preventDefault();
    setIsDragOver(fileType);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(undefined);
  };

  const handleDrop = (
    e: React.DragEvent,
    fileType: "curriculum" | "category",
  ) => {
    e.preventDefault();
    setIsDragOver(undefined);

    const files = e.dataTransfer.files;
    if (files && files[0]) {
      handleFileUpload(fileType, files[0]);
    }
  };

  return (
    <div className="flex flex-col h-full overflow-hidden">
      <div className="flex-1 overflow-y-auto">
        <div className="max-w-7xl mx-auto p-6 space-y-6">
          <Button variant="ghost" onClick={() => router.back()}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Admin
          </Button>
          <header className="mb-6">
            <h1 className="text-3xl font-bold text-chula-active">
              Add New Curriculum
            </h1>
            <p className="text-chula-active/70 mt-1">
              Upload CSV files to create a new curriculum structure
            </p>
          </header>

          <form onSubmit={handleSubmitCurriculums} className="space-y-6">
            {/* File Uploads - Curriculum & Category */}
            <Card>
              <CardHeader>
                <CardTitle>Upload Curriculum & Category Files</CardTitle>
                <CardDescription>
                  Upload the required CSV files for curriculum and categories
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {[
                    {
                      type: "curriculum" as const,
                      label: "Curriculum CSV",
                      file: form.curriculumFile,
                      templates: templates[0],
                    },
                    {
                      type: "category" as const,
                      label: "Categories CSV",
                      file: form.categoryFile,
                      templates: templates[1],
                    },
                  ].map(({ type, label, file, templates }) => (
                    <div className="space-y-2" key={type}>
                      {/* Upload Section */}
                      <div
                        className={`min-h-48 relative border-2 border-dashed rounded-xl p-8 transition-all duration-200
                      ${
                        isDragOver === type
                          ? "border-blue-400 bg-blue-50/50"
                          : "border-gray-200 hover:border-gray-300 hover:bg-gray-50/30"
                      }
                    `}
                        onDragOver={(e) => handleDragOver(e, type)}
                        onDragLeave={handleDragLeave}
                        onDrop={(e) => handleDrop(e, type)}
                      >
                        <input
                          // ref={fileInputRef}
                          id={`file-${type}`}
                          type="file"
                          accept=".csv"
                          className="hidden"
                          onChange={(e) => {
                            const selectedFile = e.target.files?.[0] || null;
                            handleFileUpload(type, selectedFile);
                          }}
                        />

                        {!file && (
                          <div className="flex flex-col items-center justify-center space-y-3">
                            <div
                              className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-100 to-indigo-100 flex items-center justify-center cursor-pointer"
                              onClick={() =>
                                document.getElementById(`file-${type}`)?.click()
                              }
                            >
                              <Upload className="w-5 h-5 text-blue-600" />
                            </div>

                            <div className="text-center">
                              <p className="text-sm font-medium text-gray-700">
                                {label}
                              </p>
                              <p className="text-sm text-gray-600 mt-1">
                                Drop your file here, or{" "}
                                <span
                                  className="text-blue-600 cursor-pointer"
                                  onClick={() =>
                                    document
                                      .getElementById(`file-${type}`)
                                      ?.click()
                                  }
                                >
                                  browse
                                </span>
                              </p>
                              <p className="text-xs text-gray-500 mt-1">
                                Supports .csv files
                              </p>
                            </div>
                          </div>
                        )}
                        {file && (
                          <div className="flex flex-col items-center justify-center space-y-3">
                            <div
                              className="w-12 h-12 rounded-full bg-gradient-to-br from-green-100 to-green-50 flex items-center justify-center cursor-pointer"
                              onClick={() =>
                                document.getElementById(`file-${type}`)?.click()
                              }
                            >
                              <FileText className="w-5 h-5 text-green-600" />
                            </div>
                            <p className="text-sm font-medium text-gray-700 text-center">
                              {file.name}
                            </p>
                          </div>
                        )}
                      </div>

                      {/* Divider */}
                      <div className="relative">
                        <div className="absolute inset-0 flex items-center">
                          <div className="w-full border-t border-gray-200" />
                        </div>
                        <div className="relative flex justify-center text-xs uppercase">
                          <span className="bg-gradient-to-br from-white to-gray-50/50 px-3 text-gray-500 font-medium">
                            or
                          </span>
                        </div>
                      </div>

                      {/* Download Section */}
                      <div>
                        <Button
                          variant="outline"
                          className="w-full"
                          onClick={() =>
                            downloadTemplate(templates.file, templates.name)
                          }
                        >
                          <Download />
                          Download Template
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Course Files - Multiple with Years */}
            <Card>
              <CardHeader>
                <CardTitle>Upload Course Files (By Year)</CardTitle>
                <CardDescription>
                  Upload course CSV files for each academic year. You can add multiple files.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Add New Course File */}
                <div className="space-y-3">
                  <div className="flex gap-3 items-end">
                    <div className="flex-1 space-y-2">
                      <label className="text-sm font-medium">Year</label>
                      <Input
                        type="text"
                        inputMode="numeric"
                        placeholder="e.g. 2566, 2567, 2568"
                        value={newCourseYear}
                        onChange={(e) => {
                          const val = e.target.value.replace(/\D/g, "");
                          setNewCourseYear(val);
                        }}
                        autoComplete="off"
                        maxLength={4}
                      />
                    </div>
                    <div className="flex-1 space-y-2">
                      <label className="text-sm font-medium">Course CSV File</label>
                      <div className="flex gap-2">
                        <Input
                          id="new-course-file"
                          type="file"
                          accept=".csv"
                          className="flex-1"
                          onChange={(e) => {
                            const file = e.target.files?.[0];
                            const year = parseInt(newCourseYear);

                            if (!file) {
                              toast.error("Please select a file");
                              return;
                            }
                            if (!year || year < 2500 || year > 2600) {
                              toast.error("Please enter a valid year (2500-2600)");
                              return;
                            }

                            handleAddCourseFile(file, year);
                            // Reset file input
                            e.target.value = "";
                          }}
                        />
                      </div>
                    </div>
                  </div>

                  {/* Download Template */}
                  <div className="flex justify-end">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        downloadTemplate(templates[2].file, templates[2].name)
                      }
                    >
                      <Download className="w-4 h-4 mr-2" />
                      Download Course Template
                    </Button>
                  </div>
                </div>

                {/* Divider */}
                {form.courseFiles.length > 0 && (
                  <div className="relative">
                    <div className="absolute inset-0 flex items-center">
                      <div className="w-full border-t border-gray-200" />
                    </div>
                    <div className="relative flex justify-center text-xs uppercase">
                      <span className="bg-white px-3 text-gray-500 font-medium">
                        Added Files
                      </span>
                    </div>
                  </div>
                )}

                {/* List of Added Course Files */}
                {form.courseFiles.length > 0 && (
                  <div className="space-y-2">
                    {form.courseFiles.map((courseFile, index) => (
                      <div
                        key={index}
                        className="flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-200"
                      >
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-green-100 to-green-50 flex items-center justify-center">
                            <FileText className="w-5 h-5 text-green-600" />
                          </div>
                          <div>
                            <p className="text-sm font-medium text-gray-900">
                              {courseFile.file.name}
                            </p>
                            <p className="text-xs text-gray-500">
                              Year: {courseFile.year}
                            </p>
                          </div>
                        </div>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => handleRemoveCourseFile(index)}
                          className="text-red-600 hover:text-red-700 hover:bg-red-50"
                        >
                          <X className="w-4 h-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                )}

                {/* Empty State */}
                {form.courseFiles.length === 0 && (
                  <div className="text-center py-8 text-gray-500 text-sm">
                    No course files added yet. Add your first course file above.
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Preview */}
            {form.curriculumFile && form.categoryFile && form.courseFiles.length > 0 && (
              <PreviewCard {...form} />
            )}

            {/* Actions */}
            <div className="flex justify-end gap-4">
              <Button variant="outline" onClick={() => router.back()}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={
                  !form.curriculumFile ||
                  !form.categoryFile ||
                  form.courseFiles.length === 0 ||
                  loading
                }
                className="bg-chula-active hover:bg-chula-active/80"
              >
                {loading ? "Creating..." : "Create Curriculum"}
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
