import 'package:rupu/domain/entities/live_voting_result.dart';

class LiveVotingMapper {
  static LiveVotingResult fromJson(Map<String, dynamic> json) {
    return LiveVotingResult(
      totalUnits: json['total_units'],
      units: List<LiveVotingUnit>.from(
        (json['units'] as List<dynamic>).map(
          (x) => LiveVotingUnit(
            propertyUnitId: x['property_unit_id'],
            propertyUnitNumber: x['property_unit_number'],
            participationCoefficient: x['participation_coefficient'],
            residentId: x['resident_id'],
            residentName: x['resident_name'],
            hasVoted: x['has_voted'],
            votingOptionId: x['voting_option_id'],
            optionText: x['option_text'],
            optionCode: x['option_code'],
            optionColor: x['option_color'],
            votedAt: x['voted_at'] != null
                ? DateTime.parse(x['voted_at'])
                : null,
          ),
        ),
      ),
    );
  }
}
