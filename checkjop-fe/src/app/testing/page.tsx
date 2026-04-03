"use client";

import { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { testCases, TestCase } from "@/data/testCases";
import { useAppStore } from "@/store/appStore";
import { Calculator, CheckCircle, XCircle, AlertCircle, Download, Upload, FileJson, PlusCircle } from "lucide-react";
import { useStudyPlan } from "@/hooks/useStudyPlan";
import { readFileAsText, parseImportFile } from "@/utils/exportImport";
import ManualCourseForm from "@/components/ManualCourseForm";

export default function TestingPage() {
  const router = useRouter();
  const { handleCalculate } = useStudyPlan();
  const {
    setStudyPlan,
    setExemptions,
    setSelectedCurriculum,
    curriculums,
    studyPlan,
    exemptions,
    selectedCurriculum,
    exportStudyPlanToJSON,
    importStudyPlanFromData,
  } = useAppStore();

  const [selectedTestCase, setSelectedTestCase] = useState<TestCase | null>(null);
  const [filter, setFilter] = useState<"all" | "implemented" | "not-implemented">("all");
  const [importStatus, setImportStatus] = useState<{
    type: "success" | "error" | "warning" | null;
    messages: string[];
  }>({ type: null, messages: [] });
  const [showManualForm, setShowManualForm] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const loadTestCase = (testCase: TestCase) => {
    // Find curriculum by name
    const curriculum = curriculums.find((c) => c.nameTH === testCase.curriculumName);

    if (!curriculum) {
      toast.error("Curriculum not found", {
        description: `Curriculum "${testCase.curriculumName}" not found. Please ensure the curriculum is loaded.`
      });
      return;
    }

    // Load the test case data into store
    setSelectedCurriculum(curriculum);
    setStudyPlan(testCase.studyPlan);
    setExemptions(testCase.exemptions);
    setSelectedTestCase(testCase);

    // Show success notification
    toast.success("Test case loaded", {
      description: `Loaded "${testCase.name}" with ${testCase.studyPlan.length} courses`
    });

    // Scroll to top after loading
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const handleCalculateAndNavigate = () => {
    if (!selectedTestCase) {
      toast.error("No test case loaded", {
        description: "Please load a test case first"
      });
      return;
    }
    handleCalculate();
    router.push("/calculate");
  };

  const handleExport = () => {
    if (studyPlan.length === 0) {
      toast.error("Cannot export", {
        description: "No study plan to export. Please add some courses first."
      });
      return;
    }
    exportStudyPlanToJSON();
    setImportStatus({
      type: "success",
      messages: ["Study plan exported successfully!"],
    });
    setTimeout(() => setImportStatus({ type: null, messages: [] }), 3000);
  };

  const handleImportClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      const content = await readFileAsText(file);
      const { data, validation } = parseImportFile(content);

      if (!validation.isValid) {
        setImportStatus({
          type: "error",
          messages: validation.errors,
        });
        return;
      }

      if (data) {
        importStudyPlanFromData(data);
        const messages = [`Successfully imported ${data.studyPlan.length} courses`];

        if (validation.warnings.length > 0) {
          setImportStatus({
            type: "warning",
            messages: [...messages, ...validation.warnings],
          });
        } else {
          setImportStatus({
            type: "success",
            messages,
          });
        }

        setSelectedTestCase(null);
        setTimeout(() => setImportStatus({ type: null, messages: [] }), 5000);
      }
    } catch (error) {
      setImportStatus({
        type: "error",
        messages: [`Failed to import file: ${(error as Error).message}`],
      });
    }

    // Reset file input
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const filteredTestCases = testCases.filter((tc) => {
    if (filter === "implemented") return tc.implemented;
    if (filter === "not-implemented") return !tc.implemented;
    return true;
  });

  return (
    <div className="container mx-auto p-6 max-w-7xl">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Test Cases Explorer</h1>
        <p className="text-gray-600">
          Load predefined test cases to test the graduation checker functionality
        </p>
      </div>

      {/* Export/Import and Manual Add Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        {/* Export/Import Section */}
        <Card className="border-purple-300 bg-purple-50/50">
          <CardHeader>
            <div className="flex items-center gap-2">
              <FileJson className="h-5 w-5 text-purple-600" />
              <CardTitle className="text-lg">Export / Import Study Plan</CardTitle>
            </div>
            <CardDescription>
              Save your current study plan to a JSON file or load a previously saved plan
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-4">
              {/* Action Buttons */}
              <div className="flex gap-3">
                <Button
                  onClick={handleExport}
                  variant="outline"
                  className="flex-1 border-green-500 text-green-700 hover:bg-green-50"
                  disabled={studyPlan.length === 0}
                >
                  <Download className="h-4 w-4 mr-2" />
                  Export to JSON
                </Button>
                <Button
                  onClick={handleImportClick}
                  variant="outline"
                  className="flex-1 border-blue-500 text-blue-700 hover:bg-blue-50"
                >
                  <Upload className="h-4 w-4 mr-2" />
                  Import from JSON
                </Button>
                <input
                  ref={fileInputRef}
                  type="file"
                  accept=".json"
                  onChange={handleFileChange}
                  className="hidden"
                />
              </div>

            {/* Current Plan Info */}
            <div className="grid grid-cols-3 gap-4 p-3 bg-white rounded-lg border text-sm">
              <div>
                <span className="font-medium text-gray-700">Courses:</span>{" "}
                <span className="text-gray-900">{studyPlan.length}</span>
              </div>
              <div>
                <span className="font-medium text-gray-700">Exemptions:</span>{" "}
                <span className="text-gray-900">{exemptions.length}</span>
              </div>
              <div>
                <span className="font-medium text-gray-700">Curriculum:</span>{" "}
                <span className="text-gray-900">{selectedCurriculum?.nameEN || "Not selected"}</span>
              </div>
            </div>

            {/* Status Messages */}
            {importStatus.type && (
              <div
                className={`p-3 rounded-lg border ${
                  importStatus.type === "success"
                    ? "bg-green-50 border-green-300 text-green-800"
                    : importStatus.type === "warning"
                    ? "bg-yellow-50 border-yellow-300 text-yellow-800"
                    : "bg-red-50 border-red-300 text-red-800"
                }`}
              >
                <div className="flex items-start gap-2">
                  {importStatus.type === "success" ? (
                    <CheckCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                  ) : importStatus.type === "warning" ? (
                    <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                  ) : (
                    <XCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                  )}
                  <div className="text-sm">
                    {importStatus.messages.map((msg, idx) => (
                      <div key={idx}>{msg}</div>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

        {/* Manual Course Add Section */}
        <Card className="border-orange-300 bg-orange-50/50">
          <CardHeader>
            <div className="flex items-center gap-2">
              <PlusCircle className="h-5 w-5 text-orange-600" />
              <CardTitle className="text-lg">Add Course Manually</CardTitle>
            </div>
            <CardDescription>
              Manually add a course to your study plan by entering its details
            </CardDescription>
          </CardHeader>
          <CardContent>
            {showManualForm ? (
              <ManualCourseForm
                onClose={() => setShowManualForm(false)}
                semester={1}
                yearOfStudy={1}
                academicYear={2566}
              />
            ) : (
              <Button
                onClick={() => setShowManualForm(true)}
                className="w-full bg-orange-600 hover:bg-orange-700"
              >
                <PlusCircle className="h-4 w-4 mr-2" />
                Open Form
              </Button>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Current Selection Info */}
      {selectedTestCase && (
        <Card className="mb-6 border-blue-500">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg">Currently Loaded: {selectedTestCase.name}</CardTitle>
                <CardDescription>{selectedTestCase.description}</CardDescription>
              </div>
              <Button onClick={handleCalculateAndNavigate} className="bg-blue-600 hover:bg-blue-700">
                <Calculator className="h-4 w-4 mr-2" />
                Calculate & View Results
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-3 gap-4 text-sm">
              <div>
                <span className="font-medium">Courses:</span> {studyPlan.length}
              </div>
              <div>
                <span className="font-medium">Exemptions:</span> {exemptions.length}
              </div>
              <div>
                <span className="font-medium">Expected:</span>{" "}
                <pre className="whitespace-pre-wrap text-xs mt-1 bg-gray-100 p-2 rounded">{selectedTestCase.expectedResult?.description || "N/A"}</pre>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Filters */}
      <div className="flex gap-2 mb-6">
        <Button
          variant={filter === "all" ? "default" : "outline"}
          onClick={() => setFilter("all")}
          size="sm"
        >
          All ({testCases.length})
        </Button>
        <Button
          variant={filter === "implemented" ? "default" : "outline"}
          onClick={() => setFilter("implemented")}
          size="sm"
        >
          Implemented ({testCases.filter((tc) => tc.implemented).length})
        </Button>
        <Button
          variant={filter === "not-implemented" ? "default" : "outline"}
          onClick={() => setFilter("not-implemented")}
          size="sm"
        >
          Not Implemented ({testCases.filter((tc) => !tc.implemented).length})
        </Button>
      </div>

      {/* Test Cases Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {filteredTestCases.map((testCase) => (
          <Card
            key={testCase.id}
            className={`cursor-pointer transition-all hover:shadow-lg ${
              selectedTestCase?.id === testCase.id ? "ring-2 ring-blue-500" : ""
            } ${!testCase.implemented ? "opacity-60" : ""}`}
            onClick={() => testCase.implemented && loadTestCase(testCase)}
          >
            <CardHeader>
              <div className="flex items-start justify-between mb-2">
                <CardTitle className="text-base">{testCase.name}</CardTitle>
                {testCase.implemented ? (
                  <Badge variant="default" className="bg-green-500">
                    <CheckCircle className="h-3 w-3 mr-1" />
                    Ready
                  </Badge>
                ) : (
                  <Badge variant="secondary">
                    <XCircle className="h-3 w-3 mr-1" />
                    Skipped
                  </Badge>
                )}
              </div>
              <CardDescription className="text-sm">{testCase.description}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2">
                  <span className="font-medium">Courses:</span>
                  <span className="text-gray-600">{testCase.studyPlan.length}</span>
                </div>
                {testCase.exemptions.length > 0 && (
                  <div className="flex items-center gap-2">
                    <span className="font-medium">Exemptions:</span>
                    <span className="text-gray-600">{testCase.exemptions.length}</span>
                  </div>
                )}
                {testCase.expectedResult && (
                  <div className="mt-3 p-2 bg-gray-50 rounded text-xs">
                    <div className="flex items-start gap-2">
                      <AlertCircle className="h-3 w-3 mt-0.5 text-blue-500 flex-shrink-0" />
                      <div>
                        <div className="font-medium mb-1">Expected Result:</div>
                        <div className="text-gray-600">
                          {testCase.expectedResult.hasViolations ? "Has violations" : "No violations"}
                        </div>
                        <pre className="text-gray-500 mt-1 whitespace-pre-wrap text-xs">{testCase.expectedResult.description}</pre>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {filteredTestCases.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          No test cases found for the selected filter
        </div>
      )}
    </div>
  );
}
