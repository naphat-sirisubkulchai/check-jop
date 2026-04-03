import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import GraphManager from "@/graph/components/graphManager"
import { CreateCurriculumForm } from "@/types"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"

const PreviewCard = (form: CreateCurriculumForm) => {
  return (
    <Card>
              <CardHeader>
                <CardTitle>Preview</CardTitle>
                <CardDescription>
                  Review the uploaded files before processing multiple
                  curriculums
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Tabs defaultValue="overview" className="w-full">
                  <TabsList className="grid w-full grid-cols-4">
                    <TabsTrigger value="overview">Overview</TabsTrigger>
                    <TabsTrigger value="details">Data Preview</TabsTrigger>
                    <TabsTrigger value="dependency">
                      Dependency Graph
                    </TabsTrigger>
                    <TabsTrigger value="files">File Info</TabsTrigger>
                  </TabsList>

                  <TabsContent value="overview" className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                        <p className="text-2xl font-bold text-gray-900">
                          {form.previewData?.totalCurriculums || 0}
                        </p>
                        <p className="text-gray-500">Curriculums</p>
                      </div>
                      <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                        <p className="text-2xl font-bold text-gray-900">
                          {form.previewData?.totalCategories || 0}
                        </p>
                        <p className="text-gray-500">Categories</p>
                      </div>
                      <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                        <p className="text-2xl font-bold text-gray-900">
                          {form.previewData?.totalCourses || 0}
                        </p>
                        <p className="text-gray-500">Courses</p>
                      </div>
                    </div>

                    {form.previewData && (
                      <div className="bg-chula-soft/20 border border-chula-active/30 rounded-lg p-4">
                        <h4 className="font-semibold text-chula-active mb-2">
                          Batch Processing Summary
                        </h4>
                        <p className="text-sm text-gray-700">
                          Ready to process {form.previewData.totalCurriculums}{" "}
                          curriculum(s) with {form.previewData.totalCategories}{" "}
                          categories and {form.previewData.totalCourses}{" "}
                          courses.
                        </p>
                      </div>
                    )}
                  </TabsContent>

                  <TabsContent value="details" className="space-y-4">
                    {form.previewData && (
                      <div className="space-y-6">
                        {/* Curriculums Preview */}
                        {form.previewData.totalCurriculums > 0 && (
                          <div>
                            <h4 className="font-semibold mb-3">
                              Curriculums ({form.previewData.totalCurriculums})
                            </h4>
                            <div className="border rounded-lg overflow-hidden">
                              <div className="max-h-40 overflow-y-auto">
                                <table className="w-full text-sm">
                                  <thead className="bg-gray-50">
                                    <tr>
                                      {form.previewData.curriculums[0] &&
                                        Object.keys(
                                          form.previewData.curriculums[0]
                                        ).map((key) => (
                                          <th
                                            key={key}
                                            className="px-3 py-2 text-left font-medium text-gray-900"
                                          >
                                            {key}
                                          </th>
                                        ))}
                                    </tr>
                                  </thead>
                                  <tbody>
                                    {form.previewData.curriculums
                                      .slice(0, 5)
                                      .map((curriculum, index) => (
                                        <tr key={index} className="border-t">
                                          {Object.values(curriculum).map(
                                            (value: any, i) => (
                                              <td key={i} className="px-3 py-2">
                                                {String(value)}
                                              </td>
                                            )
                                          )}
                                        </tr>
                                      ))}
                                  </tbody>
                                </table>
                              </div>
                              {form.previewData.totalCurriculums > 5 && (
                                <div className="bg-gray-50 px-3 py-2 text-sm text-gray-600">
                                  +{form.previewData.totalCurriculums - 5} more
                                  rows...
                                </div>
                              )}
                            </div>
                          </div>
                        )}

                        {/* Categories Preview */}
                        {form.previewData.totalCategories > 0 && (
                          <div>
                            <h4 className="font-semibold mb-3">
                              Categories ({form.previewData.totalCategories})
                            </h4>
                            <div className="border rounded-lg overflow-hidden">
                              <div className="max-h-40 overflow-y-auto">
                                <table className="w-full text-sm">
                                  <thead className="bg-gray-50">
                                    <tr>
                                      {form.previewData.categories[0] &&
                                        Object.keys(
                                          form.previewData.categories[0]
                                        ).map((key) => (
                                          <th
                                            key={key}
                                            className="px-3 py-2 text-left font-medium text-gray-900"
                                          >
                                            {key}
                                          </th>
                                        ))}
                                    </tr>
                                  </thead>
                                  <tbody>
                                    {form.previewData.categories
                                      .slice(0, 5)
                                      .map((category, index) => (
                                        <tr key={index} className="border-t">
                                          {Object.values(category).map(
                                            (value: any, i) => (
                                              <td key={i} className="px-3 py-2">
                                                {String(value)}
                                              </td>
                                            )
                                          )}
                                        </tr>
                                      ))}
                                  </tbody>
                                </table>
                              </div>
                              {form.previewData.totalCategories > 5 && (
                                <div className="bg-gray-50 px-3 py-2 text-sm text-gray-600">
                                  +{form.previewData.totalCategories - 5} more
                                  rows...
                                </div>
                              )}
                            </div>
                          </div>
                        )}

                        {/* Courses Preview */}
                        {form.previewData.totalCourses > 0 && (
                          <div>
                            <h4 className="font-semibold mb-3">
                              Courses ({form.previewData.totalCourses})
                            </h4>
                            <div className="border rounded-lg overflow-hidden">
                              <div className="max-h-40 overflow-y-auto">
                                <table className="w-full text-sm">
                                  <thead className="bg-gray-50">
                                    <tr>
                                      {form.previewData.courses[0] &&
                                        Object.keys(
                                          form.previewData.courses[0]
                                        ).map((key) => (
                                          <th
                                            key={key}
                                            className="px-3 py-2 text-left font-medium text-gray-900"
                                          >
                                            {key}
                                          </th>
                                        ))}
                                    </tr>
                                  </thead>
                                  <tbody>
                                    {form.previewData.courses
                                      .slice(0, 5)
                                      .map((course, index) => (
                                        <tr key={index} className="border-t">
                                          {Object.values(course).map(
                                            (value: any, i) => (
                                              <td key={i} className="px-3 py-2">
                                                {value}
                                              </td>
                                            )
                                          )}
                                        </tr>
                                      ))}
                                  </tbody>
                                </table>
                              </div>
                              {form.previewData.totalCourses > 5 && (
                                <div className="bg-gray-50 px-3 py-2 text-sm text-gray-600">
                                  +{form.previewData.totalCourses - 5} more
                                  rows...
                                </div>
                              )}
                            </div>
                          </div>
                        )}
                      </div>
                    )}
                  </TabsContent>

                  <TabsContent value="dependency">
                    {form.previewData && (
                      <GraphManager
                        courses={form.previewData.courses}
                        curriculums={form.previewData.curriculums}
                      />
                    )}
                  </TabsContent>

                  <TabsContent value="files" className="space-y-4">
                    <div className="space-y-3">
                      <div className="border rounded-lg p-4">
                        <div className="flex items-center justify-between">
                          <div>
                            <h4 className="font-medium text-gray-900">
                              Curriculum File
                            </h4>
                            <p className="text-sm text-gray-500">
                              {form.curriculumFile?.name}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm font-medium text-gray-900">
                              {(
                                (form.curriculumFile?.size || 0) / 1024
                              ).toFixed(1)}{" "}
                              KB
                            </p>
                            <p className="text-xs text-green-600">
                              ✓ CSV Format
                            </p>
                          </div>
                        </div>
                      </div>

                      <div className="border rounded-lg p-4">
                        <div className="flex items-center justify-between">
                          <div>
                            <h4 className="font-medium text-gray-900">
                              Categories File
                            </h4>
                            <p className="text-sm text-gray-500">
                              {form.categoryFile?.name}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm font-medium text-gray-900">
                              {((form.categoryFile?.size || 0) / 1024).toFixed(
                                1
                              )}{" "}
                              KB
                            </p>
                            <p className="text-xs text-green-600">
                              ✓ CSV Format
                            </p>
                          </div>
                        </div>
                      </div>

                      {form.courseFiles.map((courseFile, index) => (
                        <div key={index} className="border rounded-lg p-4">
                          <div className="flex items-center justify-between">
                            <div>
                              <h4 className="font-medium text-gray-900">
                                Course File (Year {courseFile.year})
                              </h4>
                              <p className="text-sm text-gray-500">
                                {courseFile.file.name}
                              </p>
                            </div>
                            <div className="text-right">
                              <p className="text-sm font-medium text-gray-900">
                                {((courseFile.file.size || 0) / 1024).toFixed(1)}{" "}
                                KB
                              </p>
                              <p className="text-xs text-green-600">
                                ✓ CSV Format
                              </p>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>

                    <div className="bg-chula-soft/20 border border-chula-active/30 rounded-lg p-4">
                      <h4 className="font-semibold text-chula-active mb-2">
                        Expected CSV Structure
                      </h4>
                      <ul className="text-sm text-gray-700 space-y-1">
                        <li>
                          • <strong>Curriculum CSV:</strong> Contains curriculum
                          information (name, year, credits, etc.)
                        </li>
                        <li>
                          • <strong>Categories CSV:</strong> Contains course
                          categories and their relationships
                        </li>
                        <li>
                          • <strong>Courses CSV (by year):</strong> Contains course
                          details and prerequisites for specific academic year
                        </li>
                      </ul>
                      <p className="text-xs text-gray-600 mt-2">
                        Note: You can upload multiple course files for different years (e.g., 2566, 2567, 2568)
                      </p>
                    </div>
                  </TabsContent>
                </Tabs>
              </CardContent>
            </Card>
  )
}

export { PreviewCard }