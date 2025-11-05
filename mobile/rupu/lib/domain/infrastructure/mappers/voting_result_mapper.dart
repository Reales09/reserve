import 'package:rupu/domain/entities/voting_results_result.dart';

class VotingResultMapper {
  static VotingResultsResult fromJson(Map<String, dynamic> json) {
    return VotingResultsResult(
      success: json['success'] ?? false,
      message: json['message'],
      data: List<VotingResult>.from(
        (json['data'] as List<dynamic>).map(
          (x) => VotingResult(
            id: x['id'],
            votingId: x['voting_id'],
            propertyUnitId: x['property_unit_id'],
            votingOptionId: x['voting_option_id'],
            votedAt: DateTime.parse(x['voted_at']),
            ipAddress: x['ip_address'],
            userAgent: x['user_agent'],
          ),
        ),
      ),
    );
  }
}
