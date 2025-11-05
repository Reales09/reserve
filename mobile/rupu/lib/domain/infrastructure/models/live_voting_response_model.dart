class LiveVotingResponseModel {
  final bool success;
  final String? message;
  final Map<String, dynamic> data;

  const LiveVotingResponseModel({
    required this.success,
    this.message,
    required this.data,
  });

  factory LiveVotingResponseModel.fromJson(Map<String, dynamic> json) {
    return LiveVotingResponseModel(
      success: json['success'] ?? false,
      message: json['message'],
      data: json['data'] as Map<String, dynamic>,
    );
  }
}
