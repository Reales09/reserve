class VotingResultsResponseModel {
  final bool success;
  final String? message;
  final List<Map<String, dynamic>> data;

  const VotingResultsResponseModel({
    required this.success,
    this.message,
    required this.data,
  });

  factory VotingResultsResponseModel.fromJson(Map<String, dynamic> json) {
    return VotingResultsResponseModel(
      success: json['success'] ?? false,
      message: json['message'],
      data: List<Map<String, dynamic>>.from(json['data'] as List<dynamic>),
    );
  }
}
